package com.vmware.suitaf.util;

import org.slf4j.LoggerFactory;

public class Logger {

   private static final org.slf4j.Logger _logger =
         LoggerFactory.getLogger(Logger.class);

   // =================================================================
   // Public logging methods
   // =================================================================
   public static void passed(String message) {
      _logger.info("PASSED: {}", message);
   }

   public static void failed(String message) {
      _logger.error("FAILED: {}", message);
   }

   public static void info(String message) {
      _logger.info(message);
   }

   public static void debug(String message) {
      _logger.info(message);
   }

   public static void warn(String message) {
      _logger.warn(message);
   }

   public static void verify(String message) {
      _logger.info("VERIFICATION: {}", message);
   }

   public static void error(String message) {
      _logger.error("Message: {}", message);
   }

   public static void error(Throwable e) {
      _logger.error("Exception: {}", e);
   }

   public static void error(String message, Throwable e) {
      _logger.error("Message: {}", message);
      _logger.error("Exception: {}", e);
   }

   public static void fatal(String message) {
      _logger.error("FATAL: {}", message);
   }

   public static void fatal(Throwable e) {
      _logger.error("FATAL: {}", e);
   }
}
