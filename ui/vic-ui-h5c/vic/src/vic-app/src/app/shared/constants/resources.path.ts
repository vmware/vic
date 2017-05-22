/*
 Copyright 2017 VMware, Inc. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

export const VIC_LOGO_100X100 = '/assets/vic-icons/100x100.png';
export const VSPHERE_VM_SUMMARY_KEY = 'vsphere.core.vm.summary';
export const VSPHERE_SERVEROBJ_VIEWEXT_KEY =
    'vsphere.core.inventory.serverObjectViewsExtension';
export const VSPHERE_VITREE_HOSTCLUSTERVIEW_KEY =
    'vsphere.core.viTree.hostsAndClustersView';
export const WS_SUMMARY = {
    keys: {
        VENDOR: 'vic_workspace.summary.vendor.label',
        VERSION: 'vic_workspace.summary.version.label',
        VCH: 'vic_workspace.summary.vch.label'
    },
    defaults: {
        'vic_workspace.summary.vendor.label': 'Vendor',
        'vic_workspace.summary.version.label': 'Version',
        'vic_workspace.summary.vch.label': null
    }
};
export const WS_VCH = {
    DG: {
        COLUMNS: {
            keys: {
                NAME: 'vic_workspace.vch.datagrid.columns.name',
                OVERALL_STATUS:
                'vic_workspace.vch.datagrid.columns.overallStatus',
                DOCKER_API_ENDPOINT:
                'vic_workspace.vch.datagrid.columns.dockerApiEndpoint',
                VCH_ADMIN_PORTAL:
                'vic_workspace.vch.datagrid.columns.vchAdminPortal'
            },
            defaults: {
                'vic_workspace.vch.datagrid.columns.name': 'Name',
                'vic_workspace.vch.datagrid.columns.overallStatus': 'Status',
                'vic_workspace.vch.datagrid.columns.dockerApiEndpoint':
                'Docker API Endpoint',
                'vic_workspace.vch.datagrid.columns.vchAdminPortal':
                'VCH Admin Portal'
            }
        }
    }
};

export const WS_CONTAINER = {
    DG: {
        COLUMNS: {
            keys: {
                CONTAINER_NAME: 'vic_workspace.container.datagrid.columns.containerName',
                POWER_STATE: 'vic_workspace.container.datagrid.columns.powerState',
                MEMORY_USAGE: 'vic_workspace.container.datagrid.columns.guestMemoryUsage',
                CPU_USAGE: 'vic_workspace.container.datagrid.columns.overallCpuUsage',
                STORAGE_USAGE: 'vic_workspace.container.datagrid.columns.committedStorage',
                PORT_MAPPING: 'vic_workspace.container.datagrid.columns.portMapping',
                VM_NAME: 'vic_workspace.container.datagrid.columns.name',
                IMAGE_NAME: 'vic_workspace.container.datagrid.columns.imageName'
            },
            defaults: {
                'vic_workspace.container.datagrid.columns.containerName': 'Name',
                'vic_workspace.container.datagrid.columns.powerState': 'State',
                'vic_workspace.container.datagrid.columns.guestMemoryUsage': 'Memory Usage',
                'vic_workspace.container.datagrid.columns.overallCpuUsage': 'CPU Usage',
                'vic_workspace.container.datagrid.columns.committedStorage': 'Storage Usage',
                'vic_workspace.container.datagrid.columns.portMapping': 'Port Mapping',
                'vic_workspace.container.datagrid.columns.name': 'VM',
                'vic_workspace.container.datagrid.columns.imageName': 'Image'
            }
        }
    }
};
