/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.diagnostics;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.BaseTestWorkflow;
import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.test.TestWorkflowStepContext;
import com.vmware.vsphere.client.automation.provider.commontb.CommonTestBedProvider;
import com.vmware.vsphere.client.automation.srv.common.spec.VcSpec;

import java.io.ByteArrayOutputStream;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;

import org.apache.http.conn.ssl.AllowAllHostnameVerifier;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.SSLContextBuilder;
import org.apache.http.conn.ssl.TrustStrategy;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.HttpResponse;
import org.testng.annotations.Test;

/**
 * Class implementing the NGC SST test scenario. The test checks that the instrumentation is enabled by
 * checking the URL https://NGC_IP/vsphere-client/diagnostic
 * NOTE: It is a temporary solution to run the ngc-sst test together with the CL run list till the respective
 * sst CAT setup is provided.
 * TODO: rkovachev to remote the test once the sst setup is provided.
 */
public class PlatformDiagnosticsTest extends BaseTestWorkflow {

   @Override
   public void initSpec() {
      BaseSpec testSpec = new BaseSpec();
      setSpec(testSpec);
      VcSpec vcToCheck = new VcSpec();
      vcToCheck.name.set(testBed.getVc());
      testSpec.add(vcToCheck);
   }

   @Override
   public void composePrereqSteps(WorkflowComposition composition) {
      // nothing to do here
   }

   @Override
   public void composeTestSteps(WorkflowComposition composition) {
      composition.appendStep(new ValidateDiagnosticsUrlLstep(),
            "Query platform diagnostics servlet");
   }

   @Override
   @Test
   @TestID(id = "0")
   public void execute() throws Exception {
      if(isProviderRun) {
         super.invokeTestExecuteCommand();
      } else {
         super.execute();
      }
   }

   @Override
   public void initSpec(WorkflowSpec testSpec, TestBedBridge testbedBridge) {
      TestbedSpecConsumer testBed = testbedBridge.requestTestbed(
            CommonTestBedProvider.class, true);

      // Spec for the VC
      VcSpec vcToCheck = testBed
            .getPublishedEntitySpec(CommonTestBedProvider.VC_ENTITY);

      testSpec.add(vcToCheck);
   }

   @Override
   public void composeTestSteps(WorkflowStepsSequence<TestWorkflowStepContext> flow) {
      flow.appendStep("Query platform diagnostics servlet", new ValidateDiagnosticsUrlLstep());
   }
}