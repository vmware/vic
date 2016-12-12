/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.suitaf.apl.webdriver;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URL;

import org.apache.http.HttpHost;
import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.message.BasicHttpEntityEnclosingRequest;
import org.jboss.netty.handler.codec.http.HttpResponseStatus;
import org.json.JSONException;
import org.json.JSONObject;
import org.openqa.selenium.remote.SessionId;

/**
 * The class contains utility methods relevant for the particular Selenium
 * WebDriver implementation
 *
 * @author itrendafilova
 */
public class WebDriverUtils {
   private static final String WEBDRIVER_SESSION_URL_FORMAT = "http://%s:%s/grid/api/testsession?session=%s";
   private static final String REQUEST_METHOD = "POST";
   private static final String NODE_URL_KEY = "proxyId";

   /**
    * Extract node ip address and port number based on a data about the hub ip
    * address, port number and session id
    *
    * @param hostName
    *           hub ip address
    * @param port
    *           hub port number
    * @param session
    *           Selenium session id
    * @return WebDriverNode
    */
   public static WebDriverNode getWebDriverNode(String hostAddress, int port,
         SessionId session) {
      WebDriverNode node = new WebDriverNode(hostAddress, port);
      try {
         HttpHost host = new HttpHost(hostAddress, port);
         HttpClient client = HttpClientBuilder.create().build();
         String sessionURLString =
               String.format(WEBDRIVER_SESSION_URL_FORMAT, hostAddress,
                     port, session);
         URL sessionURL = new URL(sessionURLString);
         BasicHttpEntityEnclosingRequest request =
               new BasicHttpEntityEnclosingRequest(REQUEST_METHOD,
                     sessionURL.toExternalForm());
         HttpResponse response = client.execute(host, request);
         int statusCode = response.getStatusLine().getStatusCode();
         if (statusCode == HttpResponseStatus.OK.getCode()) {
            JSONObject object = extractJsonObject(response);
            URL nodeURL = new URL(object.getString(NODE_URL_KEY));
            node.setHostName(nodeURL.getHost());
            node.setPort(nodeURL.getPort());
         }
      } catch (Exception e) {
         e.printStackTrace();
         throw new RuntimeException(
               "Failed to acquire remote webdriver node info. Root cause: ",
               e);
      }
      return node;
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
   private static JSONObject extractJsonObject(HttpResponse response)
         throws IOException, JSONException {
      BufferedReader rd =
            new BufferedReader(new InputStreamReader(response.getEntity()
                  .getContent()));
      StringBuffer s = new StringBuffer();
      String line;
      while ((line = rd.readLine()) != null) {
         s.append(line);
      }
      rd.close();
      JSONObject objToReturn = new JSONObject(s.toString());
      return objToReturn;
   }

}
