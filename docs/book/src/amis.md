# Pre-built Kubernetes AMIs

New AMIs are built whenever a new Kubernetes version is released for each supported OS distribution and then published to supported regions.

`clusterawsadm ami list` command lists pre-built AMIs by Kubernetes version, OS, or AWS region.
See [clusterawsadm ami list](clusterawsadm/clusterawsadm_ami_list.md) for details.

## Supported OS Distributions
- Amazon Linux 2 (amazon-2)
- Ubuntu (ubuntu-20.04, ubuntu-18.04)
- Centos (centos-7)

## Supported AWS Regions
- ap-northeast-1
- ap-northeast-2
- ap-south-1
- ap-southeast-1
- ap-northeast-2
- ca-central-1
- eu-central-1
- eu-west-1
- eu-west-2
- eu-west-3
- sa-east-1
- us-east-1
- us-east-2
- us-west-1
- us-west-2

## Most recent AMIs
<table id="amis" class="display" style="width:100%"></table>

<script>
  const amisURL = "https://d2jcv1y6kf3xwc.cloudfront.net/amis.json";
  const SEARCH_HOTKEY_KEYCODE = 83;

  // hasFocus returns true if the table search is active
  function hasFocus() {
    const tableSearchBar = document.querySelector("#amis_filter > label > input[type=search]");
    return (tableSearchBar === document.activeElement);
  }

  // Prevent the mdbook search event listener capturing the 's' key, so users can search for example 'eu-west-1
  function resetKeyHandler(e) {
    if (e.altKey || e.ctrlKey || e.metaKey || e.shiftKey || e.target.type === 'textarea' || e.target.type === 'text') { return; }

    if (e.keyCode === SEARCH_HOTKEY_KEYCODE && hasFocus()) {
        e.stopPropagation();
    }
  }

  // Insert the event listener when the document is ready
  $(function() {
    document.addEventListener('keydown', function (e) { resetKeyHandler(e); }, true);
  });

  // Table display function
  function amiListToTable(data) {
    const items = data.items.map(
      item => {

        url = `https://console.aws.amazon.com/ec2/v2/home?region=${item.spec.region}#Images:visibility=public-images;search=${item.spec.imageID};sort=name`

        imageText = `<a href="${url}">${item.spec.imageID}</a>`

        return [
          item.metadata.name,
          item.spec.os,
          item.spec.region,
          item.spec.kubernetesVersion,
          imageText,
          item.metadata.creationTimestamp,
        ]
      }
    )

    $(document).ready(function() {
      const table = $('#amis').DataTable({
        data: items,
        columns: [
          {title: "Name"},
          {title: "OS"},
          {title: "Region"},
          {title: "Kubernetes Version"},
          {title: "Image ID"},
          {title: "Creation Date"},
        ]
      })

      table
        .order([3, 'dsc'], [2, 'asc'], [1, 'asc'])
        .draw();
    });
  }


  // Lazy fetch the URL
  fetch(amisURL, {
    mode: 'cors'
  })
  .then(response => response.json())
  .then(data => amiListToTable(data))
  .catch((error) => console.error('Error:', error));
</script>

