# Getting Started
Both Developers & Operators are users of Cloudfauj.

Devs use it to deploy applications, whereas Ops use it to automate infrastructure provisioning.

Cloudfauj is a single binary that can run both server and client. Download the latest commandline app from the Github [Releases](https://github.com/cloudfauj/cloudfauj/releases) Page. Move the binary to a directory included in your system's [PATH](https://superuser.com/questions/284342/what-are-path-and-other-environment-variables-and-how-can-i-set-or-use-them).

---
**NOTE**

You can use `CGO_ENABLED=1 go build` if you want to build cloudfauj binary yourself.

---

Run `cloudfauj help` for help on all commands and options.

Cloudfauj currently only supports Mac & Linux.

## Server
Cloudfauj [Server](./concepts.md#architecture) is self-hosted. Though you can easily run it on your local workstation, we recommend running it on a VPS.

### IAM permissions
The server itself need not be hosted in AWS, but it needs credentials to be able to manage resources in your AWS account.

Since it uses the AWS Go SDK, you can supply creds using any of the supported methods in the [Credential provider chain](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials).

The following IAM policy must be granted to the server's credentials:

<details>
  <summary>Click to see the policy</summary>

  ```json
  {
      "Version": "2012-10-17",
      "Statement": [
          {
              "Sid": "CloudfaujServerPermissions",
              "Effect": "Allow",
              "Action": [
                  "ec2:*",
                  "iam:*",
                  "ecs:*"
              ],
              "Resource": ["*"]
          }
      ]
  }
  ```
</details>

### Configuration
The Server is configured via a YAML file. Below is an example, let's call it `cf-server.yml`:
```yaml
---
# Use 127.0.0.1 if you only want the server to listen on the loopback address
bind_host: '0.0.0.0'
# The TCP port for the server to listen on
bind_port: 6200
# The directory containing all internal state data of the server.
# It is very crucial that you take continuous backups of this directory.
data_dir: '/var/lib/cloudfauj'
```

### Launch
Start the server using the `server` command:

```
# As an example, we provide AWS credentials via a profile setup on the host
$ export AWS_PROFILE=my_aws_profile

$ cloudfauj server --config cf-server.yml
INFO[0000] Validating AWS credentials
INFO[0000] Setting up server data directory              dir=/var/lib/cloudfauj
INFO[0000] Downloading Terraform v1.0.4
INFO[0002] Initializing Terraform
INFO[0002] Starting CloudFauj Server                     bind_addr="0.0.0.0:6200"
```

When Cloudfauj server is started for the first time, it performs some additional tasks like setting up its base data dir and Terraform.

Try invoking the client to verify that your server is running as expected:
```
$ cloudfauj env list
No environments created yet

# Or specify the server address
$ cloudfauj --server-addr http://127.0.0.1:6200 env list
No environments created yet
```

**Previous**: [Table of Contents](../README.md#documentation)

**Next**: [Concepts](./concepts.md)
