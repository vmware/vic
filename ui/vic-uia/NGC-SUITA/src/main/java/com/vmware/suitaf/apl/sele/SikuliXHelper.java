/*
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.suitaf.apl.sele;

import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.apl.sele.sikulix.SikuliInitializer;
import org.apache.commons.collections.IteratorUtils;
import org.sikuli.basics.Settings;
import org.sikuli.script.FindFailed;
import org.sikuli.script.ImagePath;
import org.sikuli.script.Match;
import org.sikuli.script.Region;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.awt.*;
import java.util.Collections;
import java.util.Iterator;
import java.util.List;

/**
 * Wrapper class around the SikuliX library. The exposed methods should be
 * used by the SUITA framework.
 * <p/>
 * The library is initialized using the hostname of the target machine. It
 * will check whether it's a local machine or remote VNC machine and will
 * initialize the appropriate implementation.
 * <p/>
 * TODO mkalinov the class is package scoped until a better way to find the
 * image folder is found.
 */
class SikuliXHelper {
   private static final Logger _logger =
         LoggerFactory.getLogger(SikuliXHelper.class);
   private static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;
   private static final double DEFAULT_MIN_SIMILARITY = 0.7;
   private static final String SIKULI_PATH = "/sikuli";
   private final String hostName;

   /**
    * A monitor representation as a Region object. It is used to control the
    * target screen.
    *
    * NOTE: Do not use the field directly.  Instead, use the getScreen()
    * method. This is required in order to have lazy loading of the library.
    */
   private Region screen;

   public SikuliXHelper(String hostName) {
      this.hostName = hostName;
   }

   private Region getScreen() {
      if (this.screen == null) {
         screen = SikuliInitializer.initSikuliAndGetRegion(hostName);
      }
      return screen;
   }

   /**
    * Types the passed in text using the currently active Region/Screen. The
    * method expects focus to be at the desired location.  The delay
    * defines the time in milliseconds between keystrokes.
    *
    * @param text - the text to be typed out
    * @param delayInMillis - the delay in milliseconds between keystrokes
    */
   public void type(String text, int delayInMillis) {
      _logger.debug(
            String.format(
                  "Typing text: '%s', with a delay of %d milliseconds",
                  text,
                  delayInMillis));

      // Currently, passing in a delay seems to set it to 1 second
      // regardless of the input.
      //TODO mkalinov research issue and log bug to the Sikuli library
      //TODO mkalinov uncomment when the issue is fixed
      //      screen.delayType(delayInMillis);
      type(text);
   }

   /**
    * Types the passed in text using the currently active Sikul Region. The
    * method expects focus to be at the desired location.
    *
    * @param text the text to be typed.
    */
   public void type(String text) {
      getScreen().type(text);
   }

   /**
    * Finds the image on the screen and returns the middle point coordinates
    * as a Point object
    *
    * @param imageFileName the name of the image file
    *
    * @return the center point of the found match
    */
   public Point getPoint(String imageFileName) {
      Point point;
      try {
         point = getScreen().find(imageFileName).getTarget().getPoint();
      } catch (FindFailed findFailed) {
         final RuntimeException notFound = new RuntimeException(
               "Could not find the image " + "on the screen", findFailed);
         UI.audit.aplFailure(notFound);
         throw notFound;
      }

      return point;
   }

   /**
    * Finds all the instances of imageFileName on the screen and returns the
    * count
    *
    * @param imageFileName the name of the image file
    *
    * @return - found count
    */
   public int getExistingCount(String imageFileName) {
      final List list = getAll(imageFileName);
      return list.size();
   }

   /**
    * Returns all the matches found on the screen
    *
    * @param imageFileName the name of the image file
    *
    * @return {@link java.util.List} of all found matches
    */
   public List<Match> getAll(String imageFileName) {
      Iterator<Match> allResults = doGetAllMatches(imageFileName);
      return IteratorUtils.toList(allResults);
   }

   /**
    * Attempts to find the image on the screen and clicks on it.  If the
    * image cannot be found on the screen, a Runtime exception is thrown.
    * <p/>
    * The method uses the caller class, in order to find the folder where the
    * image is located.  Since SUITA has several Maven modules, it is
    * necessary to distinguish between them at runtime. The caller class
    * is used to get the absolute resource/sikuli folder where SikulX
    * should search for the image
    *
    * @param graphicalID the graphical ID of the component
    * @param callerClass the calling class
    */
   public void clickFatal(String graphicalID, Class callerClass) {
      final String callerPath = callerClass.getCanonicalName() + SIKULI_PATH;
      ImagePath.add(callerPath);
      clickFatal(graphicalID);
   }

   /**
    * Sets the accuracy with which Sikuli does pattern recognition. Valid
    * values
    * are between 0.1 (1%) and 0.99(99%). 0.99 means that Sikuli will
    * recognize only a perfect match.
    *
    * @param minSimilarity the minimum similarity which Sikuli should use
    */
   public void setRecognitionAccuracy(float minSimilarity) {
      verifyAccuracyParameter(minSimilarity);
      Settings.MinSimilarity = minSimilarity;
   }

   /**
    * Sets accuracy with which Sikuli does pattern recognition to its default
    * value of 0.7 (70%)
    */
   public void resetRecognitionAccuracy() {
      Settings.MinSimilarity = DEFAULT_MIN_SIMILARITY;
   }

   //---------------------------------------------------------------------------
   // Private methods

   /**
    * Verifies that the min similarity parameter is between 0.1 and .99
    *
    * @param minSimilarity the minimum similarity which Sikuli should use
    */
   private void verifyAccuracyParameter(float minSimilarity) {
      if (minSimilarity < 0.1 || minSimilarity > 0.99) {
         throw new IllegalArgumentException(
               "Valid values are between 0.1 and" + " 0.9");
      }
   }

   /**
    * Used by getAll to get all the matches on the screen
    *
    * @param imageFileName the name of the image file
    *
    * @return an {@link java.util.Iterator} with all found matches
    */
   private Iterator<Match> doGetAllMatches(String imageFileName) {
      Iterator<Match> allResults;
      try {
         allResults = getScreen().findAll(imageFileName);
      } catch (FindFailed findFailed) {
         _logger.warn(
               String.format(
                     "Did not find image '%s' on the " + "screen",
                     imageFileName));
         allResults = Collections.emptyIterator();
      }
      return allResults;
   }

   /**
    * Attempts to find the image on the screen and clicks on it.  If the
    * image cannot be found on the screen, a Runtime exception is thrown
    *
    * @param imageFileName - the name of the image file
    */
   private void clickFatal(final String imageFileName) {
      try {
         click(imageFileName);
      } catch (FindFailed findFailed) {
         UI.logger.error("Could not find the element on the screen.");
         throw new RuntimeException(
               "Could not find the Sikuli element on " + "the screen");
      }
   }

   /**
    * Attempts to click on the image if found on the screen.  If the image is
    * not found on the screen, a FindFailed exception is thrown.
    *
    * @param imageFileName the name of the image file
    *
    * @throws org.sikuli.script.FindFailed thrown when either not finding
    * the image file, or not finding the image on the screen.
    */
   private void click(String imageFileName) throws FindFailed {
      getScreen().click(imageFileName);
   }

   /**
    * Returns the value of a component property.
    * @param imageFileName the name of the image file
    * @param property the property to get
    * @return the value of the property
    */
   public String getSingleProperty(String imageFileName, Property property) {
      //TODO mkalinov implement
      throw new RuntimeException("Not implemented yet");
   }

   /**
    * Return the screen dimensions
    *
    * @return {@link java.awt.Dimension}
    */
   public Dimension getScreenDimensions() {
      return getScreen().getRect().getSize();
   }
}
