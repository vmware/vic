/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.datamodel;

import java.util.ArrayList;
import java.util.List;

/**
 * Represents relation between RecentTaskFilter and found in UI matching
 * RecentTask objects.
 */
public class RecentTaskFilterResult {

   /**
    * Recent tasks filter
    */
   private final RecentTaskFilter _filter;

   /**
    * List of RecentTasks that match the filter criteria
    */
   private final List<RecentTask> _matchingTasks;

   /**
    * Constructs RecentTaskFilterResults object with given filter
    *
    * @param filter
    */
   public RecentTaskFilterResult(RecentTaskFilter filter) {
      _filter = filter;
      _matchingTasks = new ArrayList<RecentTask>();
   }

   /**
    * Getter for filter
    *
    * @return filter
    */
   public RecentTaskFilter getFilter() {
      return _filter;
   }

   /**
    * Getter for matching tasks
    *
    * @return matchingTasks
    */
   public List<RecentTask> getMatchingTasks() {
      return _matchingTasks;
   }
}
