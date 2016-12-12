/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.util;

import java.util.Set;

import com.vmware.cis.authorization.client.AuthorizationServiceClient;
import com.vmware.cis.ol.client.ObjectServiceLocator;
import com.vmware.vim.binding.dataservice.ServiceInformation;
import com.vmware.vim.binding.dataservice.accesscontrol.AuthorizationService;
import com.vmware.vim.binding.dataservice.accesscontrol.internal.AuthorizationServiceInternal;
import com.vmware.vim.binding.dataservice.tagging.TagManager;
import com.vmware.vim.query.client.Client;
import com.vmware.vim.query.client.QueryAuthenticationManager;
import com.vmware.vim.query.client.QueryDispatcher;
import com.vmware.vim.query.client.exception.NotImplementedException;
import com.vmware.vim.query.client.exception.ServiceUnavailableException;

/**
 * A class that wraps the VC query client and implement its interface.
 * Used for "safe" closing the connections to the VC without explicitly calling the .close() method.
 */
public class QueryClientWrapper implements Client {
   private Client _queryClient = null;
   private boolean _closed = false;

   public QueryClientWrapper(Client queryClient) {
      _queryClient = queryClient;
   }

   @Override
   protected void finalize() {
      if (!_closed) {
         _queryClient.close();
      }
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#close()
    */
   @Override
   public void close() {
      _closed  = true;
      _queryClient.close();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#getAuthenticationManager()
    */
   @Override
   public QueryAuthenticationManager getAuthenticationManager() {
      return _queryClient.getAuthenticationManager();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#getAuthorizationService()
    */
   @Override
   public AuthorizationService getAuthorizationService() {
      return _queryClient.getAuthorizationService();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#getAuthorizationServiceInternal()
    */
   @Override
   public AuthorizationServiceInternal getAuthorizationServiceInternal() {
      return _queryClient.getAuthorizationServiceInternal();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#getObjectServiceLocator()
    */
   @Override
   public ObjectServiceLocator getObjectServiceLocator() {
      return _queryClient.getObjectServiceLocator();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#getQueryDispatcher()
    */
   @Override
   public QueryDispatcher getQueryDispatcher() {
      return _queryClient.getQueryDispatcher();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#getServiceInformation()
    */
   @Override
   public ServiceInformation getServiceInformation()
         throws ServiceUnavailableException, NotImplementedException {
      return _queryClient.getServiceInformation();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#getTagManager()
    */
   @Override
   public TagManager getTagManager() {
      return _queryClient.getTagManager();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#getVmomiClient()
    */
   @Override
   public com.vmware.vim.vmomi.client.Client getVmomiClient() {
      return _queryClient.getVmomiClient();
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#newAuthorizationServiceClient(boolean)
    */
   @Override
   public AuthorizationServiceClient newAuthorizationServiceClient(boolean arg0) {
      return _queryClient.newAuthorizationServiceClient(arg0);
   }

   /* (non-Javadoc)
    * @see com.vmware.vim.query.client.Client#newAuthorizationServiceClient(java.util.Set, java.util.Set)
    */
   @Override
   public AuthorizationServiceClient newAuthorizationServiceClient(Set<String> arg0, Set<String> arg1) {
      return _queryClient.newAuthorizationServiceClient(arg0, arg1);
   }
}
