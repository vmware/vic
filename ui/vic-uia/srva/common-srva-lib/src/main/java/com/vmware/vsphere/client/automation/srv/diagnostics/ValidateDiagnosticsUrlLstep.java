/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.diagnostics;

import java.io.ByteArrayOutputStream;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;

import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.conn.ssl.AllowAllHostnameVerifier;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.SSLContextBuilder;
import org.apache.http.conn.ssl.TrustStrategy;
import org.apache.http.impl.client.HttpClients;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

/**
 *The step verify if the build instrumented by opening the diagnostics URL.
 */
public class ValidateDiagnosticsUrlLstep  extends BaseWorkflowStep {

   private VcSpec _vcToCheck;

   @Override
   public void prepare() throws Exception {
      _vcToCheck = getSpec().get(VcSpec.class);
   }

   @Override
   public void execute() throws Exception {
      // SSLContextBuilder to plug an all-trusting TrustStrategy in
      final SSLContextBuilder builder = new SSLContextBuilder();
      builder.loadTrustMaterial(null, new TrustStrategy() {
         public boolean isTrusted(X509Certificate[] chain, String authType)
               throws CertificateException {
            return true;
         }
      });

      // SSLConnectionSocketFactory with the rigged TrustStrategy
      final SSLConnectionSocketFactory socketFactory =
            new SSLConnectionSocketFactory(builder.build(), new AllowAllHostnameVerifier());

      // A naive HttpClient :)
      final HttpClient client = HttpClients.custom()
            .setHostnameVerifier(new AllowAllHostnameVerifier())
            .setSSLSocketFactory(socketFactory)
            .build();

      final HttpResponse response = client.execute(new HttpGet("https://"
            + _vcToCheck.name.get() + "/vsphere-client/diagnostic/"));

      final ByteArrayOutputStream output = new ByteArrayOutputStream();
      response.getEntity().writeTo(output); // this clears up resources as well :)

      final String outputString = output.toString();

      if (!"Diagnostic successful.".equals(outputString)) {
         throw new Exception("Invalid diagnostic output:\\n" + outputString);
      }
   }

   @Override
   public void prepare(WorkflowSpec filteredWorkflowSpec) throws Exception {
      _vcToCheck = filteredWorkflowSpec.get(VcSpec.class);
   }

}
