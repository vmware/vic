/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.suitaf.apl.webdriver;

/**
 * Class to represent Selenium Webdriver node. The Webdriver node can be a
 * single server node or a machine from Selenium-Grid. Selenium-Grid allows you
 * run your tests on different machines against different browsers in parallel.
 */
public class WebDriverNode {

   private String hostName;
   private int port;

   public WebDriverNode(String hostName, int port) {
      super();
      this.hostName = hostName;
      this.port = port;
   }


   /**
    * Gets the host name of the Webdriver node. For a single server node this
    * method will return the server ip, for a selenium-grid the method will
    * return the IPv6 address of the Selenium node in the Selenium-Grid
    *
    * @return The host name of the Web Driver Node
    */
   public String getHostName() {
      return hostName;
   }

   /**
    * Sets the hostname of the Webdriver node
    *
    * @param hostName
    */
   public void setHostName(String hostName) {
      this.hostName = hostName;
   }

   /**
    * Gets the port number of the Web Driver Node
    *
    * @return The port number
    */
   public int getPort() {
      return port;
   }

   /**
    * Sets the port of the Webdriver node
    *
    * @param port
    */
   public void setPort(int port) {
      this.port = port;
   }

   @Override
   public String toString() {
      return "WebDriverNode [hostName=" + hostName + ", port=" + port + "]";
   }
}
