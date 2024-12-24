#!/usr/bin/env bash
# credits : https://st-g.de/2024/05/delete-unused-amis
REGION='eu-west-1'

# Initialize counters
image_count=0
total_ebs_size=0
# enable below for multi region cleanup
# for region in $(aws account list-regions --region-opt-status-contains ENABLED ENABLED_BY_DEFAULT --query "Regions[].RegionName" --output text); do
#   echo "#####################################";
#   echo "# $region";
#   echo "#####################################";
  
# Process images
for image in $(aws ec2 describe-images --owners self \
  --query 'Images[?CreationDate<`2023-01-01` && !not_null(LastLaunchedTime)] | sort_by(@, &CreationDate)[].[ImageId][]' \
  --output text --region $REGION); do
  echo "# Image ID: $image"

  # Increment image count
  ((image_count++))

  # Calculate total EBS size for the image
  ebs_size=$(aws ec2 describe-images --image-ids $image --region $REGION \
    --query 'Images[*].BlockDeviceMappings[?Ebs].Ebs.[VolumeSize][]' --output text)

  for size in $ebs_size; do
    total_ebs_size=$((total_ebs_size + size))
  done

  echo aws ec2 deregister-image --image-id $image --region $REGION; # Remove 'echo' to execute

done

# Process snapshots
for snap in $(aws ec2 describe-snapshots --owner self \
  --filters "Name=description,Values='Created by CreateImage*','Copied for DestinationAmi*'" \
  --query "Snapshots[*].SnapshotId" --output text --region $REGION); do
  echo "# Snapshot ID: $snap"
  echo aws ec2 delete-snapshot --snapshot-id $snap --region $REGION; # Remove 'echo' to execute

done

# Print summary
echo "#####################################"
echo "# Summary for region: $REGION"
echo "# Total images processed: $image_count"
echo "# Total EBS size (GiB): $total_ebs_size"
echo "#####################################"
echo "# Finished $REGION"
echo
# done
