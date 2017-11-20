# bouncer [![Download](https://api.bintray.com/packages/palantir/releases/bouncer/images/download.svg)](https://bintray.com/palantir/releases/bouncer/_latestVersion)

Bouncer rebuilds your running infrastructure to make sure it matches the infrastructure you've defined in code.

This tool inspects AWS [auto-scaling groups](https://aws.amazon.com/documentation/autoscaling/), and terminates, in a controlled fashion, any nodes whose launch configurations don't match the one currently configured on the ASG it was launched from.  It currently supports two termination methods, `serial` and `canary`; read more about them below.

Although the examples for invoking this from code below are written in Terraform, and it's convenient to be executed from within a Terraform environment, there is nothing Terraform-specific about this tool whatsoever.

## Serial

Made for bouncing nodes in ASGs of size 1 but which form a logical set.  `./bouncer serial --help` for all available options.  Ex:

```bash
./bouncer serial -a hashi-use1-stag-server-0,hashi-use1-stag-server-1,hashi-use1-stag-server-2
```

Will bounce all the servers in those ASGs one at a time, by invoking the API call `autoscaling:TerminateInstanceInAutoScalingGroup` which invokes the terminate shutdown hook (if there is one) and removes the instance from the ELB in a way that obeys the cooldown.  This API call also decrements the desired capacity by 1 when it calls this (this is to prevent a race condition from the replacement node coming up before the one to be replaced has finished being terminated), then sets it back to what you 1 when it's done.

It will also wait between each bounce for the instance to become healthy again in the ASG (so make sure you've defined either `EC2` or `ELB` health in a robust way for your use-case) before moving on to the next node in the next ASG.

### Gotchas for serial mode

This bouncer is really only intended for the ASGs its given to each be of `desired_capacity=1` (and therefore `min_size=0`).  You must set your `min_size` to at least 1 less than `desired_capacity`. If you keep `min_size`, `max_size` and `desired_capacity` all at the same value, then issuing a `TerminateInstanceInAutoScalingGroup` will actually defer the termination of the instance in question until _after_ its replacement has started and is running successfully, which is probably not what you want.

If you need to run serial mode against an ASG with an expected `desired_capacity` other than 1, you'll need to pass-in that number as part of the ASG list.  For example, if we wanted to serially bounce nomad workers.

```bash
./bouncer serial -a hashi-use1-stag-worker-linux:3,hashi-use1-stag-worker-windows:2
```

## Canary

Made for bouncing ASGs of arbitrary size where additional nodes and scale in before the old nodes have scaled out.  `./bouncer canary --help` for all available options.  Ex:

```bash
./bouncer canary -a hashi-use1-stag-worker:3
```

Only accepts 1 ASG.  In this example, the ASG must, at start time of bouncer, have `desired_capacity` of 3, and `max_size` of at least 6.  This will

* Bump desired capacity to 4 to create a "canary" node and wait for this node to become healthy.
  * If this script encounters an ASG which already has at least 1 new, healthy node, it skips the canary phase.
* Increases desired capacity to 6 to give us a whole new set of new machines.
  * If your ASG already has more than 1 healthy new node, the desired capacity is only increased to give you 3 new nodes.  That is, if bouncer started where 2 nodes were new and only 1 were old, it would skip the canary phase and instead immediately increase the `desired_capacity` to 4 to allow a 3rd new, healthy node, before terminating the lone old node.
* Terminate the old nodes decrementing the `desired_capacity` with each terminate call.
  * Bouncer calls the [ASG terminate API call](http://docs.aws.amazon.com/cli/latest/reference/autoscaling/terminate-instance-in-auto-scaling-group.html) with `should-decrement-desired-capacity`.
  * Bouncer doesn't wait for termination of nodes between terminate calls, but after all terminate calls have been issued, it waits for all nodes to finish terminating before exiting with success.

### Gotchas for canary mode

* You should have `min_size` = `desired_capacity` and have `max_size` = 2 * `desired_capacity`.
* Your ASGs [termination policy](https://www.terraform.io/docs/providers/aws/r/autoscaling_group.html#termination_policies) is ignored by this bouncer.  This bouncer instead terminates all "old" nodes one at-a-time (but without waiting for them to complete in between), instead of letting AWS do it by scaling the ASG down.  This is because the AWS scale-down will _always_ [scale down based on AZ mismatch](http://docs.aws.amazon.com/autoscaling/latest/userguide/as-instance-termination.html) before considering any other criteria.
* If any new nodes exist at bouncer start time, they must be healthy in both the EC2 and ASG, otherwise bouncer will fail immediately.
* A new node is only canaried if a healthy new node doesn't exist.  That is, if this bouncer is re-run on an ASG which already has a healthy new node, this will only add the additional new nodes to meet your final capacity without doing the entire canary workflow again (this is to reduce churn on re-runs of the script).

## Force bouncing all nodes

By default, the bouncer will ignore any nodes which are running the same launch configuration that's set on their ASG.  If you've made a change external to the launch configuration and want the bouncer to start over bouncing all nodes regardless of launch config "oldness", you can add the `-f` flag to any of the run types.  This flag marks any node whose launch time is older than the start time of the current bouncer invocation as "out of date", thus bouncing all nodes.

## Running the bouncer in Terraform

* Grab `bouncerw` at the top-level of this repo and place it in the top-level of your Terraform.
  * By top-level, this means the top-level of the whole repo, not of each module.
  * Terraform modules should _not_ need to include this wrapper, and instead each caller of modules should have one copy of the wrapper at its top-level, and all modules or top-level code should use it automatically when invoked with `./bouncerw`.
* For each logical set of ASGs you'd like cycled in this way, add a `null_resource` block to call this wrapper script and pass-in the ASG(s) in question.
  * For more information about the `null_resource` provisioner, see [the Terraform docs](https://www.terraform.io/docs/provisioners/null_resource.html).
* For example, to cycle a group of ASGs whose Terraform variable is `consul_server`, create a `null_resource` which triggers on a change to the any of the associated launch configurations:

```terraform
resource "null_resource" "consul_server_bouncer" {
  # Changes to any instance of the cluster requires re-provisioning
  triggers {
    lc_change = "${join(",", aws_autoscaling_group.consul_server.*.launch_configuration)}"
  }

  provisioner "local-exec" {
    # Redeploy all nodes in these ASGs
    command = "./bouncerw serial -a '${join(",", aws_autoscaling_group.consul_server.*.name)}'"
  }
}
```

* For an example on using bouncer in canary mode:

```terraform
resource "null_resource" "nomad_worker_bouncer" {
  # Changes to any instance of the cluster requires re-provisioning
  triggers {
    lc_change = "${aws_autoscaling_group.nomad_worker.launch_configuration}"
  }

  provisioner "local-exec" {
    # Bounce all nodes in this ASG using the canary method
    command = "./bouncerw canary -a '${aws_autoscaling_group.nomad_worker.name}:${var.worker_count}'"
  }
}
```

### Killswitch

As `bouncer` is usually going to run inside Terraform, there may be times when you want to apply all changes to your Terraform environment without actually invoking the bouncer.  In order to do this, set the `BOUNCER_KILLSWITCH` environment variable to a non-empty value.

## Running a command before instance is terminated

Sometimes, there may be an action that needs to be performed _before_ an instance is removed from its ELB.  For example, Vault listens in active/passive mode, so removing the master server from the main Vault ELB before the master has stepped-down its responsibilities to another node, means Vault is down until the lifecycle hooks kick-in and the master node steps-down.

You can use the flag `-p` to callout to an external command before every terminate call.  The number of callouts must match the number of ASGs given to bouncer (since you will probably need to pass-in the ASG name or something else to your external command).  Any waits or other checks you need before the instance is terminated should also be baked into this external command; bouncer will call terminate on the instance as soon as the external command returns success.  See below chaining example for an example of using this with Vault.

These should be used sparingly, as most logic should be baked into your AMIs terminate hook; these should only contain logic that must run before ELB removal / draining.

## Chaining bouncers together

If there are multiple ASGs in your repo which need to be bounced in a particular order, chain their associated `null_resource`s together.  Here I'm bouncing the Consul servers, then the Vault servers, and finally the Nomad workers, in order, using both forms of the bouncer.

Ex: I want all Vault nodes to be cycled after all server nodes have been cycled:

```terraform
resource "null_resource" "consul_server_bouncer" {
  # Changes to any instance of the cluster requires re-provisioning
  triggers {
    lc_change = "${join(",", aws_autoscaling_group.consul_server.*.launch_configuration)}"
  }

  provisioner "local-exec" {
    # Redeploy all nodes in these ASGs
    command = "./bouncerw serial -a '${join(",", aws_autoscaling_group.consul_server.*.name)}'"
  }
}

resource "null_resource" "vault_server_bouncer" {
  # Changes to any instance of the cluster requires re-provisioning
  triggers {
    lc_change = "${join(",", aws_autoscaling_group.vault_server_individual.*.launch_configuration)}"
  }

  provisioner "local-exec" {
    # Redeploy all nodes in these ASGs
    command = "./bouncerw serial -a '${join(",", aws_autoscaling_group.vault_server_individual.*.name)}' -p '${join(",", formatlist("./%s/vault-step-down.sh %s %s.%s", path.module, aws_autoscaling_group.vault_server_individual.*.name, var.vault_dns_name, var.zone_name))}'"
  }

  depends_on = [
    "null_resource.consul_server_bouncer",
  ]
}

resource "null_resource" "nomad_worker_bouncer" {
  # Changes to any instance of the cluster requires re-provisioning
  triggers {
    lc_change = "${aws_autoscaling_group.nomad_worker.launch_configuration}"
  }

  provisioner "local-exec" {
    # Bounce all nodes in this ASG using the canary method
    command = "./bouncerw canary -a '${aws_autoscaling_group.nomad_worker.name}:${var.worker_count}'"
  }

  depends_on = [
    "null_resource.consul_server_bouncer",
    "null_resource.vault_server_bouncer",
  ]
}
```

## Required Permissions

In order to run the bouncer, the following permissions are required:

### asg
DescribeAutoScalingGroups
DescribeLaunchConfigurations
CompleteLifecycleAction
TerminateInstanceInAutoScalingGroup
SetDesiredCapacity

### ec2
DescribeInstances
DescribeInstanceAttribute

Note that several of these permissions could cause service outages if abused.  If this is a concern, scoping the permissions is recommended.

## Contributing

For general guidelines on contributing the Palantir products, see [this page](https://github.com/palantir/gradle-baseline/blob/develop/docs/best-practices/contributing/readme.md)
