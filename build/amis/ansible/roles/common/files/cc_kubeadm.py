# Copyright 2019 The Kubernetes Authors.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
Kubeadm
-------
**Summary:** bootstrap a Kubernetes node using kubeadm

**Internal name:** ``cc_kubeadm``

**Module frequency:** per instance

**Supported distros:** all

**Config keys**::

    kubeadm:
      operation: <init/join>
      config: '/home/user/kubeadm.config'
"""

from cloudinit import util


def _run_kubeadm(config, operation):
    cmd = ["/usr/bin/kubeadm", operation, "--config", config]
    util.subp(cmd, capture=False)


def handle(name, cfg, cloud, log, _args):
    """Handler method activated by cloud-init."""

    if 'kubeadm' not in cfg:
        log.debug("Skipping 'kubeadm' module, no config")
        return

    if 'operation' not in cfg['kubeadm']:
        log.debug("Skipping 'kubeadm' module, no operation specified")
        return

    if cfg['kubeadm']['operation'] not in ('init', 'join'):
        log.debug("Skipping 'kubeadm' module, unknown operation: %s",
                  cfg['kubeadm']['operation'])
        return

    if 'config' not in cfg['kubeadm']:
        log.debug("Skipping 'kubeadm' module, no config specified")
        return

    try:
        _run_kubeadm(cfg['kubeadm']['config'], cfg['kubeadm']['operation'])
    except Exception as e:
        log.error("Failed to run kubeadm: %s" % e)

# vi: ts=4 expandtab
