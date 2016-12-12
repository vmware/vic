/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.vsphere.client.automation.host.manage.graphics.srvapi;

import static com.vmware.client.automation.common.TestSpecValidator.ensureNotNull;

import java.util.Arrays;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.host.ConfigManager;
import com.vmware.vim.binding.vim.host.GraphicsConfig;
import com.vmware.vim.binding.vim.host.GraphicsInfo;
import com.vmware.vim.binding.vim.host.GraphicsManager;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil;

import java.util.ArrayList;

public class HostGraphicsSrvApi {

   private static final String VGPU_NVIDIA_VENDOR = "NVIDIA Corporation";
   private static final String vgpuPerformanceSettings = "performance";
   private static final String vgpuSharedSettings = "shared";
   private static final Logger _logger = LoggerFactory
         .getLogger(HostGraphicsSrvApi.class);

   private static HostGraphicsSrvApi instance = null;

   protected HostGraphicsSrvApi() {
   }

   /**
    * Get instance of HostGraphicsSrvApi.
    *
    * @return created instance
    */
   public static HostGraphicsSrvApi getInstance() {
      if (instance == null) {
         synchronized (HostGraphicsSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing HostGraphicsSrvApi.");
               instance = new HostGraphicsSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Set default values - "shared" and "performance"
    *
    * @param hostSpec
    * @throws Exception
    */
   public void setHostSettings(HostSpec hostSpec) throws Exception {
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);

      ConfigManager confManager = host.getConfigManager();
      GraphicsManager graphicsManager = VcServiceUtil
            .getVcService(hostSpec.service.get())
            .getManagedObject(confManager.getGraphicsManager());

      GraphicsInfo[] graphicsInfos = graphicsManager.getGraphicsInfo();

      int graphicsInfoLength = graphicsInfos.length;
      GraphicsConfig.DeviceType[] deviceTypes = new GraphicsConfig.DeviceType[graphicsInfoLength];
      for (int i = 0; i < graphicsInfoLength; i++) {
         GraphicsConfig.DeviceType device = new GraphicsConfig.DeviceType();
         device.deviceId = graphicsInfos[i].pciId;
         device.graphicsType = vgpuSharedSettings;
         deviceTypes[i] = device;
      }

      GraphicsConfig spec = new GraphicsConfig();

      spec.hostDefaultGraphicsType = vgpuSharedSettings;
      spec.sharedPassthruAssignmentPolicy = vgpuPerformanceSettings;
      spec.deviceType = deviceTypes;
      graphicsManager.updateGraphicsConfig(spec);
   }

   /**
    * Get name of vGPU card - it is the same for all so get it from the first
    * vGPU Nvidia element
    *
    * @param hostSpec
    * @return
    * @throws Exception
    */
   public String getHostGraphicsName(HostSpec hostSpec) throws Exception {
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      String vgpuDeviceName = null;
      GraphicsInfo[] graphicsInfos = host.getConfig().getGraphicsInfo();
      for (int i = 0; i < graphicsInfos.length; i++) {
         if (graphicsInfos[i].vendorName.equals(VGPU_NVIDIA_VENDOR)
               && graphicsInfos[i].graphicsType.contains(vgpuSharedSettings)) {
            vgpuDeviceName = graphicsInfos[i].deviceName;
            break;
         }
      }
      return vgpuDeviceName;
   }

   /**
    * Return available vGPU profiles
    *
    * @param hostSpec
    * @return
    * @throws Exception
    */
   public String[] getVgpuProfiles(HostSpec hostSpec) throws Exception {
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      ConfigManager confManager = host.getConfigManager();
      GraphicsManager graphicsManager = VcServiceUtil
            .getVcService(hostSpec.service.get())
            .getManagedObject(confManager.getGraphicsManager());

      return graphicsManager.getSharedPassthruGpuTypes();
   }

   /**
    * Return pciId for VM with vGPU profile assigned
    * 
    * @param hostSpec
    * @param vmMor
    *           - VM that has assigned vGPU profile
    * @return
    * @throws Exception
    */
   public String getVmPciId(HostSpec hostSpec, ManagedObjectReference vmMor)
         throws Exception {
      ensureNotNull(hostSpec, "hostSpec is missing");
      ensureNotNull(vmMor, "vmMor is missing");
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      ConfigManager confManager = host.getConfigManager();
      GraphicsManager graphicsManager = VcServiceUtil
            .getVcService(hostSpec.service.get())
            .getManagedObject(confManager.getGraphicsManager());

      String pciId = null;
      GraphicsInfo[] graphicsInfos = graphicsManager.getGraphicsInfo();

      List<GraphicsInfo> graphicsInfosList = Arrays.asList(graphicsInfos);

      for (GraphicsInfo graphicsInfo : graphicsInfosList) {
         List<ManagedObjectReference> graphicsVms = Arrays
               .asList(graphicsInfo.getVm());
         for (ManagedObjectReference vmMorElement : graphicsVms) {
            if (vmMor.equals(vmMorElement)) {
               pciId = graphicsInfo.pciId;
               break;
            }
         }
      }
      return pciId;
   }

   /**
    * Get the number of vGPU graphics devices
    *
    * @param hostSpec
    * @return
    * @throws Exception
    */
   public Integer getVgpuDevicesNumber(HostSpec hostSpec) throws Exception {
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      return host.getConfig().getGraphicsInfo().length;
   }

   /**
    * Get pcId for all graphic devices
    *
    * @param hostSpec
    * @return
    * @throws Exception
    */
   public String[] getPciIdList(HostSpec hostSpec) throws Exception {
      HostSystem host = ManagedEntityUtil.getManagedObject(hostSpec);
      ConfigManager confManager = host.getConfigManager();
      GraphicsManager graphicsManager = VcServiceUtil
            .getVcService(hostSpec.service.get())
            .getManagedObject(confManager.getGraphicsManager());

      GraphicsInfo[] graphicsInfos = graphicsManager.getGraphicsInfo();
      List<String> graphicsPcids = new ArrayList<String>();
      List<GraphicsInfo> graphicsInfosList = Arrays.asList(graphicsInfos);
      for (GraphicsInfo graphicsInfo : graphicsInfosList) {
         graphicsPcids.add(graphicsInfo.getPciId());
      }
      String[] pcIds = graphicsPcids.toArray(new String[graphicsPcids.size()]);
      return pcIds;
   }
}