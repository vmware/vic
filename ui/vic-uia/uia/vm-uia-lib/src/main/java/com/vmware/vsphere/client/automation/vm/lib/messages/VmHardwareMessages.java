package com.vmware.vsphere.client.automation.vm.lib.messages;

import com.vmware.vsphere.client.test.i18n.gwt.Messages;

public interface VmHardwareMessages extends Messages {
   @DefaultMessage("VM Options")
   String vmOptionsTab();

   @DefaultMessage("Virtual Hardware")
   String vmVirtualHardwareTab();

   @DefaultMessage("EFI")
   String efiFirmware();

   @DefaultMessage("ESXi 6.5 and later")
   String esxCompatibility65();

   @DefaultMessage("VM version 13")
   String vmHardware13();
}
