/*

Copyright 2017 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Mac OS script starting an Ant build of the current flex project
Note: if Ant runs out of memory try defining ANT_OPTS=-Xmx512M

*/

package com.vmware.vic.mvc;

import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.HashMap;
import java.util.Map;

import javax.servlet.http.HttpServletResponse;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import com.vmware.vic.services.EchoService;


/**
 * A controller to serve HTTP JSON GET/POST requests to the endpoint "/services".
 * Its purpose is simply to redirect HTTP requests to the service APIs implemented in
 * separate components.
 */
@Controller
@RequestMapping(value = "/services")
public class ServicesController {
   private final static Log _logger = LogFactory.getLog(ServicesController.class);

   private final EchoService _echoService;

   @Autowired
   public ServicesController(
         @Qualifier("echoService") EchoService echoService) {
      _echoService = echoService;
   }

   // Empty controller to avoid compiler warnings in vic's bundle-context.xml
   // where the bean is declared
   public ServicesController() {
      _echoService = null;
   }


   /**
    * Echo a message back to the client.
    */
   @RequestMapping(value = "/echo", method = RequestMethod.POST)
   @ResponseBody
   public String echo(@RequestParam(value = "message", required = true) String message)
         throws Exception {
      return _echoService.echo(message);
   }

   /**
    * Generic handling of internal exceptions.
    * Sends a 500 server error response along with a json body with messages
    *
    * @param ex The exception that was thrown.
    * @param response
    * @return a map containing the exception message, the cause, and a stackTrace
    */
   @ExceptionHandler(Exception.class)
   @ResponseBody
   public Map<String, String> handleException(Exception ex, HttpServletResponse response) {
      response.setStatus(HttpStatus.INTERNAL_SERVER_ERROR.value());

      Map<String,String> errorMap = new HashMap<String,String>();
      errorMap.put("message", ex.getMessage());
      if(ex.getCause() != null) {
         errorMap.put("cause", ex.getCause().getMessage());
      }
      StringWriter sw = new StringWriter();
      PrintWriter pw = new PrintWriter(sw);
      ex.printStackTrace(pw);
      errorMap.put("stackTrace", sw.toString());

      return errorMap;
   }
}

