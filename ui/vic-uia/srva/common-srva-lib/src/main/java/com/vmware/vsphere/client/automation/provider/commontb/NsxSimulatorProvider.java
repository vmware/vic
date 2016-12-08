/** Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.commontb;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.StringWriter;
import java.io.Writer;
import java.net.URL;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

import org.apache.http.HttpHost;
import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.message.BasicHttpEntityEnclosingRequest;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.util.ssh.SshConnectionManager;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.commontb.spec.NsxSimulatorSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;

/**
 * This provider installs and enables the NSX simulator (nsx-transformers
 * product) on an ESX host. The simulator can only be installed with its default
 * configuration, but the ability to modify its settings will be added in the
 * near future.
 */
public class NsxSimulatorProvider implements ProviderWorkflow {

   // logger
   private static final Logger _logger = LoggerFactory
         .getLogger(NsxSimulatorProvider.class);

   public static final String DEFAULT_ENTITY = "provider.nsx.entity.default";

   // nsx-transformers build number key string in the session settings file
   private static final String TESTBED_KEY_NSX_BUILD_NUMBER = "resource.nsxsim.buildnumber";

   // URL prefix for the buildweb API
   private static final String URL_PREFIX = "http://buildapi.eng.vmware.com";
   // Path to the build deliverables
   private static final String DELIVERABLE_PREFIX = "/ob/deliverable/";
   // Template string for the request
   private static final String BUILD_PARAMS = "?build=%s&_format=json";
   // Host property for the GET request to the Buildweb API
   private static final String HTTP_HOST = "buildapi.eng.vmware.com";

   // The names of the deliverables
   private static final String NSX_ESX_DATAPATH_FILENAME = "nsx-esx-datapath.vib";
   private static final String NSX_PY_PROTOBUF_FILENAME = "nsx-py-protobuf.vib";
   private static final String NSXA_SIM_FILENAME = "nsxa-sim.zip";

   // The paths to the deliverables
   private static final String NSX_ESX_DATAPATH_PATH = "nsx-esx-datapath/esx60/nsx-esx-datapath";
   private static final String NSX_PY_PROTOBUF_PATH = "nsx-protobuf/esx55/nsx-python-protobuf";
   private static final String NSXA_SIM_PATH = "nsx-mp-agent/esx6x/nsxa-sim";

   // Esxcli presets for installing .vib files
   private static final String ESXCLI_COMMAND_VIB = "esxcli software vib install --no-sig-check -v ";
   private static final String ESXCLI_COMMAND_ZIP = "esxcli software vib install --no-sig-check -d ";

   // Temporary download directory on the ESX host
   private static final String DIRECTORY = "/tmp/nsx/";

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {
      // Request host spec for a host
      TestbedSpecConsumer hostProviderConsumer = testbedBridge
            .requestTestbed(NimbusHostProvider.class, false);
      HostSpec requestedHostSpec = hostProviderConsumer
            .getPublishedEntitySpec(NimbusHostProvider.DEFAULT_ENTITY);

      NsxSimulatorSpec nsxSimulatorSpec = new NsxSimulatorSpec();

      assemblerSpec.add(requestedHostSpec);
      assemblerSpec.add(nsxSimulatorSpec);
   }

   @Override
   public void initPublisherSpec(PublisherSpec publisherSpec) throws Exception {
      NsxSimulatorSpec nsxSpec = new NsxSimulatorSpec();
      publisherSpec.links.add(nsxSpec);
      publisherSpec.publishEntitySpec(DEFAULT_ENTITY, nsxSpec);
   }

   @Override
   public void assignTestbedSettings(PublisherSpec publisherSpec,
         SettingsReader testbedSettings) throws Exception {
      // Populate the published spec. This will be used by the
      // checkHealth and disassemble methods.
      NsxSimulatorSpec spec = (NsxSimulatorSpec) publisherSpec
            .getPublishedEntitySpec(DEFAULT_ENTITY);

      String endpoint = SettingsUtil.getRequiredValue(testbedSettings,
            BaseHostProvider.TESTBED_KEY_ENDPOINT);
      String username = SettingsUtil.getRequiredValue(testbedSettings,
            BaseHostProvider.TESTBED_KEY_USERNAME);
      String password = SettingsUtil.getRequiredValue(testbedSettings,
            BaseHostProvider.TESTBED_KEY_PASSWORD);
      String buildNum = SettingsUtil.getRequiredValue(testbedSettings,
            TESTBED_KEY_NSX_BUILD_NUMBER);

      spec.hostName.set(endpoint);
      spec.userName.set(username);
      spec.password.set(password);
      spec.buildNumber.set(buildNum);
   }

   @Override
   public void assignTestbedSettings(AssemblerSpec assemblerSpec,
         SettingsReader testbedSettings) throws Exception {
   }

   @Override
   public void assignTestbedConnectors(
         Map<ServiceSpec, TestbedConnector> serviceConnectorsMap) throws Exception {
   }

   @Override
   public void composeProviderSteps(
         WorkflowStepsSequence<? extends WorkflowStepContext> flow) throws Exception {
      flow.appendStep("Save Settings from base testbeds", new ProviderWorkflowStep() {

         private static final String REQUEST_METHOD = "GET";

         /**
          * The spec exposed by this provider
          */
         private NsxSimulatorSpec nsxSimulatorSpec;

         @Override
         public void prepare(PublisherSpec filteredPublisherSpec,
               AssemblerSpec filterAssemblerSpec, boolean isAssembling,
               SettingsReader sessionSettingsReader) throws Exception {

            if (isAssembling) {
               nsxSimulatorSpec = filterAssemblerSpec.links.get(NsxSimulatorSpec.class);
               HostSpec hostSpec = filterAssemblerSpec.links.get(HostSpec.class);

               if (nsxSimulatorSpec == null) {
                  throw new RuntimeException("Failed to populate NsxSimulatorSpec");
               }

               nsxSimulatorSpec.hostName.set(hostSpec.name.get());
               nsxSimulatorSpec.userName.set(hostSpec.userName.get());
               nsxSimulatorSpec.password.set(hostSpec.password.get());

               // format ob-xxxxxxx
               String buildNumber = SettingsUtil.getRequiredValue(sessionSettingsReader,
                     TESTBED_KEY_NSX_BUILD_NUMBER);

               if (buildNumber != null) {
                  nsxSimulatorSpec.buildNumber.set(buildNumber);
               } else {
                  _logger.error("NSX Simulator build number not found in settings file");
               }
            } else {
               // Read the published spec. Used for checkHealth and disassemble
               nsxSimulatorSpec = filteredPublisherSpec.links.get(NsxSimulatorSpec.class);
            }
         }

         @Override
         public void disassemble() throws Exception {
            SshConnectionManager sshConnectionManager = new SshConnectionManager(
                  nsxSimulatorSpec.hostName.get(), nsxSimulatorSpec.userName.get(),
                  nsxSimulatorSpec.password.get());

            uninstallVib(sshConnectionManager, NSXA_SIM_FILENAME.replaceAll(".zip", ""));
            uninstallVib(sshConnectionManager,
                  NSX_PY_PROTOBUF_FILENAME.replaceAll(".vib", ""));
            uninstallVib(sshConnectionManager,
                  NSX_ESX_DATAPATH_FILENAME.replaceAll(".vib", ""));

            sshConnectionManager.getDefaultConnection().executeSshCommand(
                  "esxcli software vib remove -n nsxa-sim", new StringWriter(),
                  new StringWriter());
         }

         @Override
         public boolean checkHealth() throws Exception {
            SshConnectionManager sshConnectionManager = new SshConnectionManager(
                  nsxSimulatorSpec.hostName.get(), nsxSimulatorSpec.userName.get(),
                  nsxSimulatorSpec.password.get());

            Writer output = new StringWriter();
            sshConnectionManager.getDefaultConnection().executeSshCommand(
                  "esxcli software vib list | grep "
                        + NSXA_SIM_FILENAME.replace(".zip", ""),
                  output, new StringWriter());

            // If the above command is successfull it will return a single line
            // 'nsxa-sim'
            return output.toString().length() > 0;
         }

         @Override
         public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {

            try {
               // Request the list of deliverables from Buildweb
               List<JSONObject> deliverables = retrieveDeliverables(
                     nsxSimulatorSpec.buildNumber.get().replaceAll("(ob-|sb-)", ""));
               // Isolate only the deliverables we need
               Map<String, String> urlMap = extractDeliverables(deliverables);

               // Connect to the host
               SshConnectionManager sshConnectionManager = new SshConnectionManager(
                     nsxSimulatorSpec.hostName.get(), nsxSimulatorSpec.userName.get(),
                     nsxSimulatorSpec.password.get());

               // prepare directory on the ESX
               sshConnectionManager.getDefaultConnection().executeSshCommand(
                     "mkdir " + DIRECTORY, new StringWriter(), new StringWriter());

               // download files on the ESX
               downloadFile(sshConnectionManager, NSX_ESX_DATAPATH_FILENAME, DIRECTORY,
                     urlMap.get(NSX_ESX_DATAPATH_FILENAME));

               downloadFile(sshConnectionManager, NSX_PY_PROTOBUF_FILENAME, DIRECTORY,
                     urlMap.get(NSX_PY_PROTOBUF_FILENAME));

               downloadFile(sshConnectionManager, NSXA_SIM_FILENAME, DIRECTORY,
                     urlMap.get(NSXA_SIM_FILENAME));

               // install simulator
               installVib(sshConnectionManager, ESXCLI_COMMAND_VIB, DIRECTORY,
                     NSX_ESX_DATAPATH_FILENAME);
               installVib(sshConnectionManager, ESXCLI_COMMAND_VIB, DIRECTORY,
                     NSX_PY_PROTOBUF_FILENAME);
               installVib(sshConnectionManager, ESXCLI_COMMAND_ZIP, DIRECTORY,
                     NSXA_SIM_FILENAME);

               sshConnectionManager.close();
               writeSettings(testbedSettingsWriter);
            } catch (JSONException e) {
               throw new RuntimeException("Error parsing JSON object: ", e);
            } catch (Exception e) {
               throw new RuntimeException("Something went wrong: ", e);
            }
         }

         /**
          * Extract json object from an http response containing the object
          *
          * @param response
          *           http response containing json object
          * @return json object wrapped in the response
          * @throws IOException
          * @throws JSONException
          */
         private JSONObject extractJsonObject(HttpResponse response)
               throws IOException, JSONException {
            BufferedReader rd = new BufferedReader(
                  new InputStreamReader(response.getEntity().getContent()));
            StringBuffer s = new StringBuffer();
            String line;
            while ((line = rd.readLine()) != null) {
               s.append(line);
            }
            rd.close();
            JSONObject objToReturn = new JSONObject(s.toString());
            return objToReturn;
         }

         /**
          * Downloads a file on the server over the provided connection
          *
          * @param connectionManager
          *           A connection manager which holds an authenticated
          *           connection to the host
          * @param directory
          *           The location of the .vib on the filesystem
          * @param filename
          *           The name of the .vib file
          * @param url
          *           The location of the .vib file on buildweb
          *
          * @throws IOException
          */
         private void downloadFile(SshConnectionManager connectionManager,
               String filename, String directory, String url) throws IOException {
            StringBuilder wgetCommandBuilder = new StringBuilder();
            wgetCommandBuilder.append("wget ");
            wgetCommandBuilder.append(" -O ");
            wgetCommandBuilder.append(directory);
            wgetCommandBuilder.append(filename);
            wgetCommandBuilder.append(" ");
            wgetCommandBuilder.append(url);

            String command = wgetCommandBuilder.toString();
            _logger.debug("Attempting to execute command: " + command);
            connectionManager.getDefaultConnection().executeSshCommand(command,
                  new StringWriter(), new StringWriter());
         }

         /**
          * Installs a .vib on the server over the provided connection
          *
          * @param connectionManager
          *           A connection manager which holds an authenticated
          *           connection to the host
          * @param cliCommand
          *           The fully constructed command for installing a .vib,
          *           excluding the path to the file itself
          * @param directory
          *           The location of the .vib on the filesystem
          * @param filename
          *           The name of the .vib file
          *
          * @throws IOException
          */
         private void installVib(SshConnectionManager connectionManager,
               String cliCommand, String directory, String filename) throws IOException {
            StringBuilder cliCommandBuilder = new StringBuilder();
            cliCommandBuilder.append(cliCommand);
            cliCommandBuilder.append(directory);
            cliCommandBuilder.append(filename);

            String command = cliCommandBuilder.toString();
            _logger.debug("Attempting to execute command: " + command);
            connectionManager.getDefaultConnection().executeSshCommand(command,
                  new StringWriter(), new StringWriter());
         }

         /**
          * Retrieves the information for the deliverables of the specified
          * build
          *
          * @param buildNumber
          *           The build number of the nsx simulator.
          *
          * @throws Exception
          */
         private List<JSONObject> retrieveDeliverables(String buildNumber)
               throws Exception {
            // There are a little over 330 deliverables of the nsx sim build
            List<JSONObject> deliverables = new ArrayList<>(350);

            HttpHost host = new HttpHost(HTTP_HOST);
            HttpClient client = HttpClientBuilder.create().build();
            URL url = new URL(String
                  .format(URL_PREFIX + DELIVERABLE_PREFIX + BUILD_PARAMS, buildNumber));
            BasicHttpEntityEnclosingRequest request = new BasicHttpEntityEnclosingRequest(
                  REQUEST_METHOD, url.toExternalForm());

            // get all pages
            while (true) {
               HttpResponse response = client.execute(host, request);
               // Status code 200 - OK
               if (response.getStatusLine().getStatusCode() == 200) {
                  JSONObject object = extractJsonObject(response);
                  JSONArray list = object.getJSONArray("_list");
                  for (int i = 0; i < list.length(); i++) {
                     deliverables.add(list.getJSONObject(i));
                  }
                  Object nextUrl = object.get("_next_url");
                  if (JSONObject.NULL.equals(nextUrl)) {
                     break;
                  }
                  url = new URL(URL_PREFIX + nextUrl);
                  request = new BasicHttpEntityEnclosingRequest(REQUEST_METHOD,
                        url.toExternalForm());
               } else {
                  break;
               }
            }

            return deliverables;
         }

         /**
          * Maps and returns the required .vib/.zip resources in the following
          * format: (filename, download_url)
          *
          * @param deliverables
          *           The list of all available deliverables for the build
          * @return
          * @throws JSONException
          */
         private Map<String, String> extractDeliverables(List<JSONObject> deliverables)
               throws JSONException {
            Map<String, String> urlMap = new HashMap<>();

            // The JSON property for the download url
            String downloadUrlProp = "_download_url";

            Iterator<JSONObject> iterator = deliverables.iterator();
            while (iterator.hasNext()) {
               JSONObject current = iterator.next();
               String path = current.getString("path");

               if (path.contains(NSX_ESX_DATAPATH_PATH) && path.contains(".vib")) {
                  urlMap.put(NSX_ESX_DATAPATH_FILENAME,
                        current.getString(downloadUrlProp));
               } else if (path.contains(NSX_PY_PROTOBUF_PATH) && path.contains(".vib")) {
                  urlMap.put(NSX_PY_PROTOBUF_FILENAME,
                        current.getString(downloadUrlProp));
               } else if (path.contains(NSXA_SIM_PATH) && path.contains(".zip")) {
                  urlMap.put(NSXA_SIM_FILENAME, current.getString(downloadUrlProp));
               }
            }

            return urlMap;
         }

         /**
          * Write to the settings file
          */
         private void writeSettings(SettingsWriter writer) {
            writer.setSetting(BaseHostProvider.TESTBED_KEY_ENDPOINT,
                  nsxSimulatorSpec.hostName.get());
            writer.setSetting(BaseHostProvider.TESTBED_KEY_USERNAME,
                  nsxSimulatorSpec.userName.get());
            writer.setSetting(BaseHostProvider.TESTBED_KEY_PASSWORD,
                  nsxSimulatorSpec.password.get());
            writer.setSetting(TESTBED_KEY_NSX_BUILD_NUMBER,
                  nsxSimulatorSpec.buildNumber.get());
         }

         /**
          * Uninstalls the vib with the specified name from the ESX host
          * 
          * @throws IOException
          */
         private void uninstallVib(SshConnectionManager connectionManager, String name)
               throws IOException {
            connectionManager.getDefaultConnection().executeSshCommand(
                  "esxcli software vib remove -n " + name, new StringWriter(),
                  new StringWriter());
         }
      });
   }

   @Override
   public int providerWeight() {
      // TODO Auto-generated method stub
      return 0;
   }

   @Override
   public Class<? extends ProviderWorkflow> getProviderBaseType() {
      // TODO Auto-generated method stub
      return null;
   }
}
