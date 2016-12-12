/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

public class NicSpec extends ManagedEntitySpec {

   /**
    * Address type values
    *
    * <ol>
    * <li><code>Manual</code> - Statically assigned MAC address.</li>
    * <li><code>Generated</code> - Automatically generated MAC address.</li>
    * <li><code>Assigned</code> - MAC address assigned by VirtualCenter.</li>
    * </ol>
    */
   public enum AddressType {
      MANUAL("Manual"), GENERATED("Generated"), ASSIGNED("Assigned");

      private String value;

      private AddressType(String value) {
         this.value = value;
      }

      public String value() {
         return value;
      }
   }

   /**
    * The {@code AdapterType} values are used to indicate the supported device types
    * for network adapter.
    */
   public enum AdapterType {
      E1000("E1000"), E1000E("E1000e"), VMXNET3("VMXNET 3");

      private String value;

      private AdapterType(String value) {
         this.value = value;
      }

      public String value() {
         return this.value;
      }
   }

   /**
    * Indicates whether the device is currently connected.
    * Valid only while the virtual machine is running.
    */
   public DataProperty<Boolean> connected;

   /**
    * Specifies whether or not to connect the device when the virtual machine starts.
    */
   public DataProperty<Boolean> startConnected;

   /**
    * Specifies the address type.
    */
   public DataProperty<AddressType> addressType;

   /**
    * The device type of the network adapter.
    */
   public DataProperty<AdapterType> adapterType;

   /**
    * The name of the device on the host system.
    */
   public DataProperty<String> deviceName;

   /**
    * MAC address assigned to the virtual network adapter. Set this property only if
    * <code>addressType</code> is of type <i>Manual</i>
    */
   public DataProperty<String> macAddress;

   /**
    * Indicates whether the nic will be added or edited
    */
   public enum NicOperationType {
      EDIT("Edit"), ADD("Add");

      private String value;

      private NicOperationType(String value) {
         this.value = value;
      }

      public String value() {
         return value;
      }
   }

   /**
    * Specifies if nic will be added or edited
    */
   public DataProperty<NicOperationType> nicOperationType;

   /**
    * 0-based index of the network adapter. This index is required when
    * {@code nicOperationType} is {@link NicOperationType.EDIT}.
    */
   public DataProperty<Integer> networkAdapterIndex;

   /**
    * Specifies the dvportgroup of the nic
    */
   public DataProperty<DvPortgroupSpec> dvPortGroup;

   /**
    * Specifies the opaque network of the nic
    */
   public DataProperty<OpaqueNetworkSpec> opaqueNetwork;

   /**
    * The resource allocation settings of the nic
    */
   public DataProperty<NetworkResourceAllocationSpec> resourceAllocation;
}
