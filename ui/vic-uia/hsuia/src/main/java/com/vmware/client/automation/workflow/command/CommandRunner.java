/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.command;

import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.hsua.common.util.IOUtils;

/**
 * Class used to run workflow commands. It will be used by test systems to run
 * providers and tests without using the TesNG.
 * The class executes the commands provided by the input parameter.
 */
public class CommandRunner {

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(CommandRunner.class);

   /**
    * The method expects absolute path to the file with list of commands to
    * execute.
    * @param args
    * @throws Exception
    */
   public static void main(String[] args) throws Exception {

      if(args.length == 0) {
         _logger.info("No command file is provided to the Command Runner!");
         return;
      }
      String commandsFilePath =  args[0];;
      List<String> commandsList = IOUtils.readLinesFromFile(commandsFilePath);

      CommandController controller = new CommandController();
      controller.initialize(commandsList);
      controller.prepare();
      controller.execute();
   }
}
