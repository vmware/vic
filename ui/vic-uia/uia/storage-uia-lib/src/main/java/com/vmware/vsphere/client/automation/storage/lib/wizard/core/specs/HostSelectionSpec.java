package com.vmware.vsphere.client.automation.storage.lib.wizard.core.specs;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;

/**
 * {@link EntitySpec} implmenetation adding verification data for the state of
 * the host which can be selected
 */
public class HostSelectionSpec extends EntitySpec {

   /**
    * The spec describing the host
    */
   public DataProperty<HostSpec> hostSpec;

   /**
    * Should the host be enabled
    */
   public DataProperty<Boolean> isEnabledExpected;

   /**
    * The suffix to the host display name
    */
   public DataProperty<String> hostDisplayNameSuffix;

}
