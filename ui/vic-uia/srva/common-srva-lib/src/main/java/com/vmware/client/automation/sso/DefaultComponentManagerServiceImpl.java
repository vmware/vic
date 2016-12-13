/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.sso;

import org.apache.commons.lang.NotImplementedException;

import com.vmware.vim.binding.cis.cm.SearchCriteria;
import com.vmware.vim.binding.cis.cm.ServiceInfo;
import com.vmware.vim.binding.cis.cm.fault.ComponentManagerFault;
import com.vmware.vise.vim.cm.ComponentManagerService;

/**
 * A very basic implementation of {@link ComponentManagerService} interface.
 * Receives component manager URL and {@link ServiceInfo} object for the SSO service
 * in the constructor and implements only the methods that return them.
 */
class DefaultComponentManagerServiceImpl implements ComponentManagerService {

   private final String _cmServiceUrl;
   private final ServiceInfo _ssoServiceInfo;

   /**
    * Creates implementation of ComponentManagerService and wraps the given
    * component manager URL and {@code ServiceInfo} object for SSO service.
    *
    * @param cmServiceUrl component manager URL
    * @param ssoServiceInfo {@code ServiceInfo} object for SSO service
    */
   DefaultComponentManagerServiceImpl(String cmServiceUrl, ServiceInfo ssoServiceInfo) {
      _cmServiceUrl = cmServiceUrl;
      _ssoServiceInfo = ssoServiceInfo;
   }

   @Override
   public String getServiceUrl() {
      return _cmServiceUrl;
   }

   @Override
   public ServiceInfo[] searchSso() throws ComponentManagerFault {
      return new ServiceInfo[] { _ssoServiceInfo };
   }

   @Override
   public ServiceInfo[] search(SearchCriteria searchCriteria)
         throws ComponentManagerFault {
      throw new NotImplementedException();
   }

   @Override
   public ServiceInfo getService(String s) throws ComponentManagerFault {
      throw new NotImplementedException();
   }

   @Override
   public String getHostId() {
      throw new NotImplementedException();
   }
}
