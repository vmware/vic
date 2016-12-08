package com.vmware.hsua.common.datamodel;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;


/**
 * This class implements a sparse graph container with directed and undirected
 * edges and typed vertices.<br>
 * (See: <a href=http://en.wikipedia.org/wiki/Graph_theory>Graph Theory</a>)
 * <br>
 * It will be used to store the links between {@link PropertyBox} instances. The
 * meaning of each link will be determined by the specific types of the
 * two instances. Each of the links will be either directed or undirected.
 * <br>
 * @author dkozhuharov
 */
public final class SparseGraph {

    /**
     * This enumeration represents the types of links that could be established
     * between a couple of instances in a {@link SparseGraph}. Each type of
     * link has its <b>"pair"</b> type. Here is the list of pairs:
     * <li> {@link #UNDIRECTED} has pair type {@link #UNDIRECTED}
     * <li> {@link #OUTGOING} has pair type {@link #INCOMING}
     * <li> {@link #INCOMING} has pair type {@link #OUTGOING}
     * <br><br>
     * When a link is established with one of the link types - automatically
     * a reverse link is established with the <b>"pair"</b> type. For example:
     * <br>
     * Adding link: A {@link #UNDIRECTED} B
     * creates automatically: B {@link #UNDIRECTED} A
     * <br>
     * Adding link: A {@link #OUTGOING} B
     * creates automatically: B {@link #INCOMING} A
     * <br>
     * Adding link: A {@link #INCOMING} B
     * creates automatically: B {@link #OUTGOING} A
     * <br><br>
     * @author dkozhuharov
     */
    public static enum LinkType {
        /**
         * This link type represents "undirected" link. For two linked
         * instances <b>A</b> and <b>B</b> we have the following to be true:
         * <li> <b>A</b> is linked to <b>B</b>
         * <li> <b>B</b> is linked to <b>A</b>
         */
        UNDIRECTED,
        /**
         * This link type represents "outgoing" link. For two linked
         * instances <b>A</b> and <b>B</b> we have the following to be true:
         * <li> <b>A</b> is <b><i>source</i></b> for <b>B</b>
         * <li> <b>B</b> is <b><i>target</i></b> for <b>A</b>
         */
        OUTGOING,
        /**
         * This link type represents "incoming" link. For two linked
         * instances <b>A</b> and <b>B</b> we have the following to be true:
         * <li> <b>A</b> is <b><i>target</i></b> for <b>B</b>
         * <li> <b>B</b> is <b><i>source</i></b> for <b>A</b>
         */
        INCOMING;

        private LinkType pairType = null;
        static {
            LinkType.UNDIRECTED.pairType = LinkType.UNDIRECTED;
            LinkType.OUTGOING.pairType = LinkType.INCOMING;
            LinkType.INCOMING.pairType = LinkType.OUTGOING;
        }
    }

    @SuppressWarnings("serial")
    private class NeighbourMap extends HashMap<Class<?>, List<Object>> {};

    /**
     * The keys of this map are the Vertexes of the LinksGraph. The value
     * corresponding to each Vertex is another map - a neighbor map.
     * The neighbor map of a vertex contains lists of its LinksGraph neighbors.
     * The lists are split by the class of their elements. The keys of
     * the neighbor map are class instances and their corresponding list
     * contains only Vertex neighbors of this class.
     */
    private final HashMap<LinkType, HashMap<Object, NeighbourMap>> links;

    /**
     * This default constructor prepares the base structures that will hold
     * the graph linking information.
     */
    SparseGraph() {
        links = new HashMap<LinkType, HashMap<Object, NeighbourMap>>();
        links.put(LinkType.UNDIRECTED, new HashMap<Object, NeighbourMap>());
        links.put(LinkType.INCOMING, new HashMap<Object, NeighbourMap>());
        links.put(LinkType.OUTGOING, new HashMap<Object, NeighbourMap>());
    }

    private final void addHalf(
            LinkType lType, Object linkPeer1, Object linkPeer2) {
        if (linkPeer1 == null || linkPeer2 == null) {
        	throw new IllegalArgumentException("Linked Graph nodes must not be nulls.");
        }

        HashMap<Object, NeighbourMap> components = links.get(lType);

        if (!components.containsKey(linkPeer1)) {
            components.put(linkPeer1, new NeighbourMap());
        }
        NeighbourMap neighbours = components.get(linkPeer1);

        if (!neighbours.containsKey(linkPeer2.getClass())) {
            neighbours.put(linkPeer2.getClass(), new ArrayList<Object>());
        }
        List<Object> typedNeighbours = neighbours.get(linkPeer2.getClass());

        if (!typedNeighbours.contains(linkPeer2)) {
            typedNeighbours.add(linkPeer2);
        }
    }
    /**
     * This method allows adding a single link between two instances passed
     * with the parameters <b>linkPeer1</b> and <b>linkPeer2</b>. The created
     * link will be of {@link LinkType} given by parameter <b>lType</b>.
     * Automatic link in the opposite direction is created with the
     * <b>"pair"</b> link type (see details in {@link LinkType}).
     * @param lType - the type of the link to be created
     * @param linkPeer1 - the first instance to be linked
     * @param linkPeer2 - the second instance to be linked
     */
    public final void addLink(
            LinkType lType, Object linkPeer1, Object linkPeer2) {
        addHalf(lType, linkPeer1, linkPeer2);
        addHalf(lType.pairType, linkPeer2, linkPeer1);
    }

    private final void removeHalf(
            LinkType lType, Object linkPeer1, Object linkPeer2) {
        NeighbourMap neighbours = links.get(lType).get(linkPeer1);
        if (neighbours == null) {
            return;
        }

        List<Object> typedNeighbours = neighbours.get(linkPeer2.getClass());
        if (typedNeighbours == null) {
            return;
        }

        typedNeighbours.remove(linkPeer2);
        if (typedNeighbours.size() == 0) {
            neighbours.remove(linkPeer2.getClass());
        }
        if (neighbours.size() == 0) {
            links.get(lType).remove(linkPeer1);
        }
    }
    /**
     * This method allows removing a single link between two instances passed
     * with the parameters <b>linkPeer1</b> and <b>linkPeer2</b>. The removed
     * link will be of {@link LinkType} given by parameter <b>lType</b>.
     * Automatically the link in the opposite direction is removed. It must be
     * of the <b>"pair"</b> link type (see details in {@link LinkType}).
     * <br><br>
     * <b>NOTE:</b> No exception is thrown if no such link exists.
     * <br><br>
     * @param lType - the type of the link to be removed
     * @param linkPeer1 - the first instance of the removed link
     * @param linkPeer2 - the second instance of the removed link
     */
    public final void removeLink(
            LinkType lType, Object linkPeer1, Object linkPeer2) {
        removeHalf(lType, linkPeer1, linkPeer2);
        removeHalf(lType.pairType, linkPeer2, linkPeer1);
    }

    /**
     * This method retrieves all instances that are linked from <b>linkPeer1</b>
     * instance and satisfy the following conditions:
     * <li> Are instances of type given by <b>linkPeer2Class</b>
     * <li> Are linked through a link of type given by <b>lType</b>
     * <br><br>
     * @param lType - the type of the link to be searched
     * @param linkPeer1 - the instance from which the link must start
     * @param linkPeer2Class - the class of the linked instances that must be
     * retrieved
     * @return a list of instances compliant to the search criteria
     */
    public final List<Object> getLinks(
            LinkType lType, Object linkPeer1, Class<?> linkPeer2Class) {
        List<Object> linksOfType = new ArrayList<Object>();

        NeighbourMap neighbours = links.get(lType).get(linkPeer1);

        if (neighbours != null) {
            for (Class<?> neighboursClass : neighbours.keySet()) {
                if (linkPeer2Class.isAssignableFrom(neighboursClass)) {
                    linksOfType.addAll(neighbours.get(neighboursClass));
                }
            }
        }

        return linksOfType;
    }

   @Override
   public String toString() {
      return links.toString();
   }
}

