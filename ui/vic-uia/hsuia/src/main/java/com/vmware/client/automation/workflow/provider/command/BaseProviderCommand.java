/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider.command;

import com.vmware.client.automation.workflow.command.BaseCommand;

/**
 * Holds
 *
 */
public abstract class BaseProviderCommand extends BaseCommand {

   /**
    * Provider commands
    * =====================================
    * 
    * assemble-testbed:
    * - [In] Canonical provider class names (Ex: a.b.c.common-testbed-provider)
    * - [Out] Test bed instance properties - key-value file (Ex: end-point IP settings, user, pass, other testbed artifacts)
    * 
    * disassemble-testbed:
    * - [In] Canonical provider class names (Ex: a.b.c.common-testbed-provider)
    * - [In] Test bed instance properties - key-value file (Ex: end-point IP settings, user, pass, other testbed artifacts)
    * 
    * check-testbed:
    * - [In] Canonical provider class names (Ex: a.b.c.common-testbed-provider)
    * - [In] Test bed instance properties - key-value file (Ex: end-point IP settings, user, pass, other testbed artifacts)
    * 
    * list-providers:
    * - [In] Package name (Ex: a.b.c)
    * - [Out] Canonical provider class names - CVS file
    * 
    * map-provider:
    * - Canonical provider name (Ex: a.b.c.common-testbed-provider)
    * - Output CVS file
    *    -- Example:
    *          a.b.c.common-testbed-provider, a.b.c.cloudvm-provider
    *          a.b.c.common-testbed-provider, a.b.c.host-provider
    *          a.b.c.common-testbed-provider, a.b.c.host-provider
    * 
    */

   /**
    * Registry commands (TODO: Define output
    * =====================================
    * 
    * init-registry:
    *  - provider properties fileS
    *  - context properties file -> here come all type of key values such as the log file, screenshot folder and many other.
    *  - Log file
    *  - Screenshot output folder?
    * 
    *  registry-init
    *   - [In] Global Settings File
    * 
    *  registry-set-testbed
    *   - [In] Testbed settings file
    */

   /**
    * Test commands
    * =====================================
    * 
    * test-list:
    * - [In] Package name
    * - [Out] Canonical test class names - CVS file
    * 
    * test-analyze:
    * - [In] Canonical test class name
    * - [Out] Required providers - CVS file
    * - [Out] Test steps - CVS file
    * 
    * test-execute:
    * - [In] Canonical test class name
    * - [In] Test scope
    * - [Out] Report - plain text file
    */
}
