/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */

package com.vmware.utils.ssl;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLException;
import javax.net.ssl.SSLSession;
import java.security.cert.Certificate;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;

public class ThumbprintHostNameVerifier implements HostnameVerifier {

   @Override
   public boolean verify(String host, SSLSession session) {
      try {
         Certificate[] certificates = session.getPeerCertificates();
         verify(host, (X509Certificate) certificates[0]);
         return true;
      } catch (SSLException e) {
         return false;
      }
   }

   private void verify(String host, X509Certificate cert) throws SSLException {
      try {
         String thumbprint = ThumbprintTrustManager.getThumbprint(cert);
         ThumbprintTrustManager.checkThumbprint(thumbprint);
      } catch(CertificateException e){
         throw new SSLException(e.getMessage());
      }
   }
}
