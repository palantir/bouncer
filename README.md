# bouncer

Bouncer rebuilds your running infrastructure to make sure it matches the infrastructure you've defined in code. All releases can be downloaded from [the release page](https://github.com/palantir/bouncer/releases), or automated in your code via bouncerw (see below).

This tool inspects AWS [auto-scaling groups](https://aws.amazon.com/documentation/autoscaling/), and terminates, in a controlled fashion, any nodes whose launch templates or launch configurations don't match the one currently configured on the ASG it was launched from.  It currently supports two termination methods, `serial` and `canary`; read more about them below.

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

## Rolling

Rolling is the same as serial, but does not decrement. This means `min_size`, `max_size`, and `desired_capacity` can all be the same value.

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

## Slow-canary

Similar to canary, but instead of adding a single canary host, waiting for it to become healthy, then adding the remainder, it only adds one new node at a time, destroying an old node as it goes. Can be useful to avoid quorum issues if halving a cluster all at once is too much churn.

Another way to think about it would be it's like serial mode for ASGs with more than 1 instance, where we don't ever want to drop below _desired capacity_ number of nodes in service.

Ex:

```bash
./bouncer slow-canary -a hashi-use1-stag-server:3
```

Only accepts 1 ASG. In this example, the ASG must, at start time of bouncer, have `desired_capacity` of 3, and `max_size` of at least 4.  This will

* Bump desired capacity to 4 to create our first "canary" node and wait for this node to become healthy.
  * So that we will keep the number of active nodes at any given time to either 3 or 4, and not let it get to 2 or 5.
* Terminates one of the old nodes, NOT decrementing the `desired_capacity` along with it, so that AWS can replace it with a fresh node.
  * Bouncer calls the [ASG terminate API call](http://docs.aws.amazon.com/cli/latest/reference/autoscaling/terminate-instance-in-auto-scaling-group.html) setting `should-decrement-desired-capacity` to false.
* Once nodes have settled and we have again 4 healthy nodes, we call terminate on another old node.
* Once nodes have settled and we have again 4 healthy nodes (3 being new), we call terminate WITH `should-decrement-desired-capacity` to let us go back to our steady-state of 3 nodes.

## New experimental batch modes

Eventually, `batch-serial` and `batch-canary` could potentially replace `serial` and `canary` with their default values. However, given that the logic in the new batch modes is significantly different than the older modes, they're implemented in parallel for now to prevent potential disruptions relying on the existing behaviour.

### Batch-canary

Use-case is, you'd prefer to use canary, but you don't have capacity (or money) to double your ASG in the middle phase. Instead, you want to add new nodes in batches as you delete old nodes.

This method takes in a `batch` parameter. Your final desired capacity + `batch` determines the maximum bouncer will ever scale your ASG to. Bouncer will never scale your ASG below your desired capacity.

I.e. the core tenants of batch-canary are:

* Your total `InService` node count will not go below your given desired capacity
* Your total node count regardless of status will never go above (desired capacity + batch size)

NOTE: You should probably suspend the "AZ Rebalance" process on your ASG so that AWS doesn't violate these contraints either.

EX: You have an ASG of size 4. You don't have enough instance capacity to run 8 instances, but you do have enough to run 6. Invoke bouncer in `batch-canary` with a `batchsize` of `2` to accomplish this. This will

* Bump desired capacity to 5 to create our first "canary" node.
* Wait for this node to become healthy.
* Kill an old node. Note that we will wait for it to leave `Terminating`, and in this case also wait for `Terminating:Wait` to complete before proceeding. This is because we will wait for the ASG to "settle" to the desired capacity chosen.
* Given we just killed a node, desired capacity is back to `4`. Canary phase is done.
* Set desired capacity to our max size, `6`, which starts two new machines spawning.
* Wait for the highly ephemeral phase of `Pending` to complete, but do NOT wait for `Pending:Wait` to complete.
* Kill 2 old nodes to get us back down to `4` desired capacity. We've now issued kills to 3 old nodes.
* Again wait for the ASG to settle, which means waiting for the nodes we just killed to fully terminate.
* Given we only have one old node now, we increase our desired capacity only up to 5, to give us one new node.
* Once this node leaves `Pending` and enters `Pending:Wait`, issue a terminate to the final old node.
* Wait for the old node to totally finish terminating, AND wait for all in-flight new nodes (in this case 1) to become `InService` before completing

### Batch-serial

Use-case is, you'd prefer to use serial, but you have way too many instances so this takes too long. You can't use canary because your desired capacity is also your max capacity for external reasons (perhaps you tie a static EBS volume to every instance in this ASG).

This method takes in a `batch` parameter. Your final desired capacity - `batch` determines the maximum number of instances bouncer will delete at any time. Bouncer will never scale your ASG above your desired capacity.

I.e. the core tenants of batch-serial are:

* Your total `InService` node count will not go below (desired capacity - batch size)
* Your total node count regardless of status will never go above your given desired capacity

NOTE: You should probably suspend the "AZ Rebalance" process on your ASG so that AWS doesn't violate these contraints either.

EX: You have an ASG of size 4. You don't want to delete one instance at a time, but two at a time is ok. Set `batch` to `2`. This mode still does canary a single node so you don't potentially batch a huge number of instances that might all fail to boot. This will

* Terminate a single node, waiting for it to be fully destroyed.
* Increase desired capacity back to original value.
* Wait for this new node to become healthy.
* Kill up to batchsize nodes, so in this case, 2. Wait for them to fully die.
* Increase desired capacity back to original value, and wait for all nodes to come up healthy.
* Kill last old node, wait for it to fully die.
* Increase capacity back to original value, and wait for all nodes to become healthy.

## Force bouncing all nodes

By default, the bouncer will ignore any nodes which are running the same launch template version (or same launch configuration) that's set on their ASG.  If you've made a change external to the launch configuration / template and want the bouncer to start over bouncing all nodes regardless of launch config / template "oldness", you can add the `-f` flag to any of the run types.  This flag marks any node whose launch time is older than the start time of the current bouncer invocation as "out of date", thus bouncing all nodes.

## Running the bouncer in Terraform

* Grab `bouncerw` at the top-level of this repo and place it in the top-level of your Terraform.
  * By top-level, this means the top-level of the whole repo, not of each module.
  * Terraform modules should _not_ need to include this wrapper, and instead each caller of modules should have one copy of the wrapper at its top-level, and all modules or top-level code should use it automatically when invoked with `./bouncerw`.
* For each logical set of ASGs you'd like cycled in this way, add a `null_resource` block to call this wrapper script and pass-in the ASG(s) in question.
  * For more information about the `null_resource` provisioner, see [the Terraform docs](https://www.terraform.io/docs/provisioners/null_resource.html).
* For example, if you're using launch tamplates, to cycle an ASG in canary mode, create a `null_resource` which triggers on a change to the launch template:

```terraform
resource "null_resource" "server_canary_bouncer" {
  triggers {
    trigger = "${aws_launch_template.server.latest_version}"
  }

  provisioner "local-exec" {
    command = "./bouncerw canary -a '${aws_autoscaling_group.server.name}:${var.worker_count}'"
  }
}
```

* If you need to use serial mode across multiple ASGs, something like this can work:

```terraform
resource "null_resource" "server_serial_bouncer" {
  triggers {
    trigger = "${join(",", aws_launch_template.server.*.latest_version)}"
  }

  provisioner "local-exec" {
    command = "./bouncerw serial -a '${join(",", aws_autoscaling_group.server.*.name)}'"
  }
}
```

* If you're using launch configs instead, it's similar, but the trigger needs to be the LC:

```terraform
resource "null_resource" "server_bouncer" {
  triggers {
    trigger = "${aws_autoscaling_group.server.launch_configuration}"
  }

  provisioner "local-exec" {
    command = "./bouncerw canary -a '${aws_autoscaling_group.server.name}:${var.server_count}'"
  }
}
```

* And of course the similar example in serial mode w/ multiple ASGs:

```terraform
resource "null_resource" "server_canary_bouncer" {
  triggers {
    trigger = "${join(",", aws_autoscaling_group.server.*.launch_configuration)}"
  }

  provisioner "local-exec" {
    command = "./bouncerw serial -a '${join(",", aws_autoscaling_group.server.*.name)}'"
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

If there are multiple ASGs in your repo which need to be bounced in a particular order, chain their associated `null_resource`s together.  Here I'm bouncing the Consul servers, then the Vault servers, then Nomad servers, and finally the Nomad workers, in order.

```terraform
resource "null_resource" "consul_server_bouncer" {
  triggers {
    trigger = "..."
  }

  provisioner "local-exec" {
    command = "..."
  }
}

resource "null_resource" "vault_server_bouncer" {
  triggers {
    trigger = "..."
  }

  provisioner "local-exec" {
    command = "..."
  }

  depends_on = [
    "null_resource.consul_server_bouncer",
  ]
}

resource "null_resource" "nomad_server_bouncer" {
  triggers {
    trigger = "..."
  }

  provisioner "local-exec" {
    command = "..."
  }

  depends_on = [
    "null_resource.consul_server_bouncer",
    "null_resource.vault_server_bouncer",
  ]
}

resource "null_resource" "nomad_worker_bouncer" {
  triggers {
    trigger = "..."
  }

  provisioner "local-exec" {
    command = "..."
  }

  depends_on = [
    "null_resource.consul_server_bouncer",
    "null_resource.nomad_server_bouncer",
    "null_resource.vault_server_bouncer",
  ]
}
```

## Required Permissions

In order to run the bouncer with launch templates, the following permissions are required:

```
autoscaling:DescribeAutoScalingGroups
autoscaling:CompleteLifecycleAction
autoscaling:TerminateInstanceInAutoScalingGroup
autoscaling:SetDesiredCapacity
ec2:DescribeInstances
ec2:DescribeInstanceAttribute
ec2:DescribeLaunchTemplates
```

For using bouncer with launch configurations, the required permissions are:

```
autoscaling:DescribeAutoScalingGroups
autoscaling:DescribeLaunchConfigurations
autoscaling:CompleteLifecycleAction
autoscaling:TerminateInstanceInAutoScalingGroup
autoscaling:SetDesiredCapacity
ec2:DescribeInstances
ec2:DescribeInstanceAttribute
```

Note that several of these permissions could cause service outages if abused.  If this is a concern, scoping the permissions is recommended.

## Contributing

For general guidelines on contributing the Palantir products, see [this page](https://github.com/palantir/gradle-baseline/blob/develop/docs/best-practices/contributing/readme.md)
