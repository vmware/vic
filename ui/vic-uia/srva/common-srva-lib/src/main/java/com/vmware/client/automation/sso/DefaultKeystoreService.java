/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.sso;

import java.security.KeyStore;
import java.security.cert.X509Certificate;

import com.vmware.client.automation.exception.SsoException;
import com.vmware.vim.binding.cis.cm.ServiceEndPoint;
import com.vmware.vim.binding.cis.cm.ServiceEndPointType;
import com.vmware.vim.binding.cis.cm.ServiceInfo;
import com.vmware.vise.util.security.CertificateUtil;
import com.vmware.vise.vim.commons.ssl.KeystoreService;
import com.vmware.vise.vim.security.sso.SsoConstants;

/**
 * A very basic implementation of {@link KeystoreService} interface. Receives
 * {@link ServiceInfo} object for the SSO service and obtains from it the SSL
 * certificates and {@link KeyStore}.
 */
public class DefaultKeystoreService implements KeystoreService {

   private X509Certificate[] _sslTrustCertificates;
   private KeyStore _keyStore;

   /**
    * Constructor.
    * 
    * Obtains from given {@link ServiceInfo} object for the SSO service and the
    * SSL certificates and {@link KeyStore}.
    * 
    * @param ssoServiceInfo
    *           {@link ServiceInfo} object for the SSO service
    */
   public DefaultKeystoreService(ServiceInfo ssoServiceInfo)
         throws SsoException {
      try {
         ServiceEndPoint ssoAdminService = getSsoAdminEndpoint(ssoServiceInfo);

         if (ssoAdminService != null) {
            String trustedAnchor = ssoAdminService.getSslTrust()[0];
            X509Certificate cert = (X509Certificate) CertificateUtil
                  .generateCertificate(trustedAnchor);

            _sslTrustCertificates = new X509Certificate[] { cert };
            _keyStore = CertificateUtil.getKeyStore(cert);
         }
      } catch (Exception e) {
         throw new SsoException(e);
      }
   }

   @Override
   public KeyStore getKeyStore() {
      return _keyStore;
   }

   @Override
   public String getKeyStorePassword() {
      return null;
   }

   @Override
   public X509Certificate[] getSslTrustCertificates() {
      return _sslTrustCertificates;
   }

   private ServiceEndPoint getSsoAdminEndpoint(ServiceInfo service) {
      for (ServiceEndPoint endPoint : service.getServiceEndPoints()) {
         ServiceEndPointType type = endPoint.getEndPointType();
         if (SsoConstants.SSO_ADMIN_ENDPOINT_TYPE.equals(type.getTypeId())) {
            return endPoint;
         }
      }
      return null;
   }
}
