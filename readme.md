# AWS commands

```
# see instances
aws ec2 describe-instances --profile costamagna-terraform --region us-east-1 --query "Reservations[*].Instances[*].[Tags[?Key=='Name']|[0].Value, InstanceId, InstanceType, State.Name, PublicIpAddress]" --output table 

# start instance 
aws ec2 start-instances --profile costamagna-terraform --region us-east-1 --instance-ids {instance-id}

# stop instance
aws ec2 stop-instances --profile costamagna-terraform --region us-east-1 --instance-ids {instance-id}

```