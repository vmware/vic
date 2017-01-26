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

package com.vmware.vic.services;

/**
 * Service handling some actions invoked from the UI
 *
 * It must be declared as osgi:service with the same name in
 * main/resources/META-INF/spring/bundle-context-osgi.xml
 */
public interface SampleActionService {
   /**
    * Sample action called on the server.
    *
    * @param objRef   Internal reference to the vCenter object for that action.
    */
   public void sampleAction1(Object objRef);

   /**
    * Sample action called on the server.
    *
    * @param objRef   Internal reference to the vCenter object for that action.
    * @return true is the action is successful, false otherwise.
    */
   public boolean sampleAction2(Object objRef);
}
