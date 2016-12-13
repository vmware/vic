/*
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.suitaf.apl.sele.sikulix;

import edu.unh.iol.dlc.ConnectionController;
import edu.unh.iol.dlc.VNCScreen;
import org.sikuli.script.Region;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.Socket;

/**
 * Implementation of the AbstractSikuliScreenFactory which creates and
 * returns remote VNC Sikuli Screen/Region objects.
 * <p/>
 * Since the SikuliX VNC code is still a little rough around the edges, the
 * process of creating a new VNC connection involves the following steps:<br>
 * 1. Open a Socket connection to the target machine<br>
 * 2. Create a ConnectionController and set it up (The connection controller
 * uses the socket connection)<br>
 * 3. Create a VNC Screen (Creation of the screen depends on the successful
 * creation of a Socket and ConnectionController instances)
 */
class RemoteVncRegionFactory extends AbstractRegionFactory {
   private static final Logger _logger = LoggerFactory.getLogger(
         RemoteVncRegionFactory.class);
   /**
    * Socket timeout in milliseconds
    */
   private static final int SOCKET_TIMEOUT = 2000;

   /**
    * VNC connection pixel format.  Truecolor and Colormap are currently
    * supported
    */
   private static final String PIXEL_FORMAT = "Truecolor";

   /**
    * VNC connection bits per pixel. Truecolor valid values are 8,16 or 32.
    */
   private static final int BITS_PER_PIXEL = 32;

   /**
    * VNC connection big endian flag (0 is little endian).
    */
   private static final int BIG_ENDIAN_FLAG = 0;

   /**
    * IP of the target node where VNC is running
    */
   private final String vncIp;
   /**
    * Port to the remote VNC service
    */
   private final Integer vncPort;

   /**
    * Socket object used by Sikuli to communicate with the remote VNC host
    */
   private Socket socket;

   /**
    * A {@link org.sikuli.script.Region} object used to control the remote
    * display.
    */
   private Region screen;

   /**
    * The ConnectionController class manages all of the VNC connections as well
    * as the local copies of the remote Framebuffers.  A thread (VNCThread) is
    * created to manage the data from each connection.  Connection
    * Controller also extends GraphicsEnvironment so that it can be used
    * with the Java2D API
    */
   private ConnectionController controller;

   /**
    * Index of the first connection thread. Used to refer to the
    * only connection we are establishing to the remote VNC machine.
    */
   private static final int FIRST_CONNECTION_THREAD = 0;

   /**
    * Creates a RemoteVncRegionFactory using the target machine's IP and port
    * . A VNC server should be running on the target and using the port.
    * <p/>
    * The factory can establish and return {@link org.sikuli.script.Region}
    * objects, via which the remote screen can be controlled.
    *
    * @param vncIp the IP of the remote VNC enabled machine
    * @param vncPort the VNC port
    */
   public RemoteVncRegionFactory(String vncIp, String vncPort) {
      verifyParameters(vncIp, vncPort);
      this.vncIp = vncIp;
      this.vncPort = Integer.valueOf(vncPort);
   }

   @Override
   public Region getRegion() {
      openSocketToVncMachine();
      createAndSetupConnectionController();
      createVNCScreen();
      return screen;
   }

   //---------------------------------------------------------------------------
   // Private methods

   /**
    * Verifies that the parameters are valid - non null.
    *
    * @param vncIp
    * @param vncPort
    */
   private void verifyParameters(String vncIp, String vncPort) {
      if (vncIp == null || vncIp.isEmpty()) {
         throw new IllegalArgumentException("VNC IP is null");
      }
      if (vncPort == null || vncPort.isEmpty()) {
         throw new IllegalArgumentException("VNC port is null");
      }
   }

   /**
    * Opens a socket connection to the VNC machine. Throws a RuntimeException
    * if the connection was not opened.
    */
   private void openSocketToVncMachine() {
      _logger.debug(String.format("Opening socket to %s:%d", vncIp, vncPort));
      try {
         socket = new Socket(vncIp, vncPort);
         socket.setSoTimeout(SOCKET_TIMEOUT);
         socket.setKeepAlive(true);
      } catch (IOException ex) {
         throw new RuntimeException("Could not create socket.", ex);
      }
   }

   /**
    * Creates and setups the options for a Sikuli ConnectionController object.
    */
   private void createAndSetupConnectionController() {
      _logger.debug("Creating and setting up ConnectionController.");
      controller = new ConnectionController(socket);
      openConnectionToRemoteMachine();
      setupPixelFormat();
      establishConnection();

      //TODO mkalinov check that the connection is really OK
   }

   /**
    * Starts the first connection thread on the ConnectionController object.
    */
   private void establishConnection() {
      controller.start(FIRST_CONNECTION_THREAD);
   }

   /**
    * Sets up the pixel format which will be used by the VNC communication.
    */
   private void setupPixelFormat() {
      controller.setPixelFormat(
            FIRST_CONNECTION_THREAD,//index of the connection thread
            PIXEL_FORMAT,//we want Truecolor representation of the remote screen
            BITS_PER_PIXEL,//we want 32bits per pixel in our representation
            BIG_ENDIAN_FLAG);//Big endian flag used in VNC communication
   }

   private void openConnectionToRemoteMachine() {
      controller.openConnection(FIRST_CONNECTION_THREAD);
   }

   /**
    * Creates a new VNCScreen
    */
   private void createVNCScreen() {
      _logger.debug("Creating VNC screen.");
      screen = new VNCScreen();
   }
}
