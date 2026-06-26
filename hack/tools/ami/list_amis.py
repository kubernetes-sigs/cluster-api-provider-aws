#!/usr/bin/env python3
"""List AMIs whose kubernetes_version tag contains a given string, across multiple regions."""

import argparse
import sys
import boto3
from botocore.exceptions import ClientError


def parse_args():
    parser = argparse.ArgumentParser(
        description="List AMIs filtered by kubernetes_version tag substring."
    )
    parser.add_argument(
        "--regions",
        default=",".join([
            "ap-south-1", "eu-west-3", "eu-west-2", "eu-west-1",
            "ap-northeast-2", "ap-northeast-1", "sa-east-1", "ca-central-1",
            "ap-southeast-1", "ap-southeast-2", "eu-central-1",
            "us-east-1", "us-east-2", "us-west-1", "us-west-2",
        ]),
        help="Comma-separated AWS regions to scan (default: all major regions)",
    )
    parser.add_argument(
        "--version-contains",
        required=True,
        metavar="STRING",
        help="Substring to match against the kubernetes_version tag value",
    )
    parser.add_argument(
        "--owner",
        default="self",
        help="AMI owner filter (default: self)",
    )
    return parser.parse_args()


def find_amis(region: str, version_contains: str, owner: str) -> list:
    """Return AMIs in region whose kubernetes_version tag value contains version_contains."""
    ec2 = boto3.client("ec2", region_name=region)
    paginator = ec2.get_paginator("describe_images")
    pages = paginator.paginate(
        Owners=[owner],
        Filters=[
            {
                "Name": "tag-key",
                "Values": ["kubernetes_version"],
            }
        ],
    )
    matches = []
    for page in pages:
        for image in page["Images"]:
            for tag in image.get("Tags", []):
                if tag["Key"] == "kubernetes_version" and version_contains in tag["Value"]:
                    matches.append(image)
                    break
    return matches


def main():
    args = parse_args()
    regions = [r.strip() for r in args.regions.split(",")]
    total = 0
    region_counts = {}

    print(f"Filter: kubernetes_version contains '{args.version_contains}'")
    print(f"Regions: {regions}")

    for region in regions:
        print(f"\n--- {region} ---")
        try:
            amis = find_amis(region, args.version_contains, args.owner)
        except ClientError as e:
            print(f"  ERROR: {e}", file=sys.stderr)
            continue

        if not amis:
            print("  No matching AMIs found.")
            region_counts[region] = 0
            continue

        for ami in amis:
            kv = next(
                t["Value"] for t in ami["Tags"] if t["Key"] == "kubernetes_version"
            )
            print(f"  {ami['ImageId']}  kubernetes_version={kv}  name={ami.get('Name', '')}")

        count = len(amis)
        print(f"  Count: {count}")
        region_counts[region] = count
        total += count

    print(f"\n--- Summary ---")
    for region, count in region_counts.items():
        print(f"  {region}: {count}")
    print(f"  Total: {total}")


if __name__ == "__main__":
    main()
