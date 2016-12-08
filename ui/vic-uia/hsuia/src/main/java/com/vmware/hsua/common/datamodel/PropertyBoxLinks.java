package com.vmware.hsua.common.datamodel;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;

import com.vmware.hsua.common.datamodel.SparseGraph.LinkType;

/**
 * This class contains the mechanics that allow creation and maintenance of
 * links between {@link PropertyBox} instances (<b>PBoxes</b>).
 * For convenience each {@link PropertyBoxLinks} instance wraps one <b>PBox</b>
 * instanced called further <b>Managed PBox</b>.<br>
 * This class contains accessor methods that enable setting, retrieval and
 * removal of links between the <b>Managed PBox</b> and other <b>PBox</b>
 * instances.<br>
 * The links information is stored in a static structure declared within this
 * class. In this structure all links are associated with the {@link Thread}
 * that is adds the linking information.<br>
 * Also there is a special "cleaner" thread, which checks and releases all
 * link information of the closed threads. This garbage-collection mechanism
 * is aligned to the fact that the <b>TestNG</b> framework executes each test
 * method in a new separate thread.
 * <br><br>
 * <b>NOTE:</b> Linking <b>PBoxes</b> is very important feature because it
 * enables creation of arbitrary data structures. <b>PBox</b> structures are
 * needed in UI testing to represent the <i>UI state</i> or the <i>back-end
 * state</i> of the tested application. Also these data structures are needed
 * to represent <i>planned</i>, <i>present</i> or <i>expected</i> states
 * needed for the implementation of each test case.<br>
 * <br>
 * @author dkozhuharov
 */
public class PropertyBoxLinks {
   private final PropertyBox self;
   /**
    * This constructor accepts a reference to the {@link PropertyBox} instance
    * whose links will be managed by this {@link PropertyBoxLinks} instance.
    * @param self - the <b>Managed PBox</b> instance
    */
   PropertyBoxLinks(PropertyBox self) {
      this.self = self;
   }

   private static final LinkType DEFAULT_LTYPE = LinkType.OUTGOING;
   private final LinkType addDefault(LinkType lType) {
      return (lType == null)? DEFAULT_LTYPE: lType;
   }
   private final LinkType[] addDefault(LinkType[] lTypes) {
      return (lTypes == null || lTypes.length == 0)?
            new LinkType[] {DEFAULT_LTYPE}: lTypes;
   }

   /**
    * This method retrieves <b>one</b> {@link PropertyBox} instance,
    * which has been linked to the <b>Managed PBox</b>.
    * The returned instance must comply with two conditions:
    * <li> Be an instance of class given by <b>linkedPBoxClass</b>
    * <li> The link to them must be of type given by <b>linkTypes</b>
    * <br><br>
    * <b>NOTE:</b> If no {@link LinkType} parameter is given the default is
    * {@link LinkType#UNDIRECTED}.
    * <br>
    * <b>NOTE:</b> If more than one instance is found - only the first one is
    * returned. No exception is thrown in such case.
    * <br><br>
    * @param <L> - the base type of the linked instances that are searched
    * @param linkedPBoxClass - a {@link Class} instance that represents the
    * base type of the linked instances that are searched
    * @param linkTypes - (Optional) the type of the links that are searched
    * (see: {@link LinkType}).
    * @return a {@link PropertyBox} instance that meets the search criterion
    * or <b>null</b> if none is found
    * @see #getAll(Class, LinkType...)
    */
   public <L extends PropertyBox> L get(
         Class<L> linkedPBoxClass, LinkType...linkTypes) {
      List<L> all = getAll(linkedPBoxClass, linkTypes);
      return all.size()>0? all.get(0): null;
   }
   /**
    * This method retrieves <b>all</b> {@link PropertyBox} instances,
    * which has been linked to the <b>Managed PBox</b>.
    * The returned instances must comply with two conditions:
    * <li> Be an instance of class given by <b>linkedPBoxClass</b>
    * <li> The link to them must be of type given by <b>linkTypes</b>
    * <br><br>
    * <b>NOTE:</b> If no {@link LinkType} parameter is given the default is
    * {@link LinkType#UNDIRECTED}.
    * <br><br>
    * @param <L> - the base type of the linked instances that are searched
    * @param linkedPBoxClass - a {@link Class} instance that represents the
    * base type of the linked instances that are searched
    * @param linkTypes - (Optional) the type of the links that are searched
    * (see: {@link LinkType}).
    * @return a list of {@link PropertyBox} instances that meets the search
    * criterion or empty list if none is found
    */
   @SuppressWarnings("unchecked")
   public <L extends PropertyBox> List<L> getAll(
         Class<L> linkedPBoxClass, LinkType...linkTypes) {
      List<L> all = new ArrayList<L>();
      for (LinkType linkType : addDefault(linkTypes)) {
         List<Object> tmp = getPBLinksGraph().getLinks(linkType, self, linkedPBoxClass);
         for (Object object : tmp) {
            all.add( (L)object );
         }
      }
      return all;
   }

   /**
    * This method adds a link connecting the <b>Managed PBox</b> to one
    * or more {@link PropertyBox} instances given as enumerated list. The
    * type of the links is the default type {@link LinkType#UNDIRECTED}.
    * <br>
    * @param linkedPBoxes - vararg that could take one or more
    * {@link PropertyBox} instances for linking.
    */
   public void add(PropertyBox...linkedPBoxes) {
      add(null, linkedPBoxes);
   }
   /**
    * This method adds a link connecting the <b>Managed PBox</b> to one
    * or more {@link PropertyBox} instances given in iterable container. The
    * type of the links is the default type {@link LinkType#UNDIRECTED}.
    * <br>
    * @param linkedPBoxes - iterable container that holds a list of
    * {@link PropertyBox} instances for linking.
    */
   public void add(Iterable<? extends PropertyBox> linkedPBoxes) {
      add(null, linkedPBoxes);
   }
   /**
    * This method adds a link connecting the <b>Managed PBox</b> to one
    * or more {@link PropertyBox} instances given as enumerated list. The
    * type of the links is given by the parameter <b>linkType</b>. If omitted
    * the type is set by default to {@link LinkType#UNDIRECTED}.
    * <br>
    * @param linkType - the type of the links to be added.
    * @param linkedPBoxes - vararg that could take one or more
    * {@link PropertyBox} instances for linking.
    */
   public void add(LinkType linkType, PropertyBox...linkedPBoxes) {
      add(linkType, Arrays.asList(linkedPBoxes));
   }
   /**
    * This method adds a link connecting the <b>Managed PBox</b> to one
    * or more {@link PropertyBox} instances given in iterable container. The
    * type of the links is given by the parameter <b>linkType</b>. If omitted
    * the type is set by default to {@link LinkType#UNDIRECTED}.
    * @param linkType - the type of the links to be added.
    * @param linkedPBoxes - vararg that could take one or more
    * {@link PropertyBox} instances for linking.
    */
   public void add(
         LinkType linkType, Iterable<? extends PropertyBox> linkedPBoxes) {
      for (PropertyBox linkPeer2 : linkedPBoxes) {
         getPBLinksGraph().addLink(addDefault(linkType), self, linkPeer2);
      }
   }

   /**
    * This method removes a link between the <b>Managed PBox</b> and one
    * {@link PropertyBox} instances. The link must be of the type given by
    * parameter <b>linkTypes</b>. If omitted the type is set by default to
    * {@link LinkType#UNDIRECTED}.
    * @param linkedPBox - the {@link PropertyBox} instances to be disconnected
    * @param linkTypes - (Optional) the type of the link to be removed
    */
   public void remove(PropertyBox linkedPBox, LinkType...linkTypes) {
      for (LinkType linkType : addDefault(linkTypes)) {
         getPBLinksGraph().removeLink(linkType, self, linkedPBox);
      }
   }
   /**
    * This method removes the links between the <b>Managed PBox</b> and all
    * {@link PropertyBox} instances of the class given by parameter
    * <b>linkedPBoxClass</b>. The links must be of the type given by
    * parameter <b>linkTypes</b>. If omitted the type is set by default to
    * {@link LinkType#UNDIRECTED}.
    * @param linkedPBoxClass - the class of the {@link PropertyBox} instances
    * to be disconnected
    * @param linkTypes - (Optional) the type of the links to be removed
    */
   public void remove(
         Class<? extends PropertyBox> linkedPBoxClass, LinkType...linkTypes){
      for (PropertyBox linkedPBox : getAll(linkedPBoxClass, linkTypes)) {
         remove(linkedPBox, linkTypes);
      }
   }

   // ======================================================================
   // Property Box Links-Graph register
   // ======================================================================

   /**
    * {@link SparseGraph} instance register
    */
   static final HashMap<Thread, SparseGraph> LinksGraphRegister =
         new HashMap<Thread, SparseGraph>();
   static final Object LinksGraphRegisterLock =
         new Object();

   /**
    * This method allows synchronized access to the {@link SparseGraph}
    * instance of the current thread. If such instance was not yet created -
    * creates it and returns it back.
    * <br>
    * @return a {@link SparseGraph} instance
    */
   static final SparseGraph getPBLinksGraph() {
      Thread currThread = Thread.currentThread();
      SparseGraph pbLinksGraph = null;

      synchronized (LinksGraphRegisterLock) {
         if (!LinksGraphRegister.containsKey(currThread)) {
            LinksGraphRegister.put(currThread, new SparseGraph());

            // Wake up the cleaner to check for old abandoned LinksGraphs
            LinksGraphCleaner.interrupt();
         }
         pbLinksGraph = LinksGraphRegister.get(currThread);
      }

      return pbLinksGraph;
   }

   // ======================================================================
   // Property Box Links-Graph cleaner
   // ======================================================================

   /**
    * This field holds a reference to a separate daemon {@link Thread}
    * process. It runs the data cleaning logic, that takes care of the
    * {@link SparseGraph} instances whose Thread objects are non-active.
    */
   static final Thread LinksGraphCleaner =
         new Thread(null, null, "LinksGraph Cleaner Thread") {
      @Override
      public void run() {
         // The cleaners endless loop
         while (true) {
            try {
               // a delay between regular checks
               sleep(1000);
            } catch (InterruptedException e) {
            }

            // call the cleanup routine
            finalize();
         }
      }

      @Override
      public void finalize() {
         synchronized (LinksGraphRegisterLock) {
            for (Thread t :
               LinksGraphRegister.keySet().toArray(new Thread[0])) {
               if (!t.isAlive()) {
                  LinksGraphRegister.remove(t);
               }
            }
         }
      }
   };

   static {
      LinksGraphCleaner.setDaemon(true);
      LinksGraphCleaner.start();
   }

   @Override
   public String toString() {
      SparseGraph graph = getPBLinksGraph();
      return graph != null ? graph.toString() : super.toString();
   }

}