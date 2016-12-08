package com.vmware.suitaf.apl.sele;

import com.google.common.base.Preconditions;
import com.vmware.suitaf.apl.ComponentMatcher;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.IDPart;

public class SeleComponentID implements ComponentMatcher {
    public final DirectID mainID;
    public final DirectID parentID;

    public SeleComponentID(SeleAPLImpl apl, IDGroup...idGroup) {
        DirectID mainID = null;
        DirectID parentID = null;

        for (IDGroup group : idGroup) {
            if (group == null) {
                continue;
            }

            DirectID id = DirectID.from(apl, group);

            if (group.has(IDPart.GROUPROLE.MAIN)) {
                mainID = id;
            }
            if (group.has(IDPart.GROUPROLE.PARENT)) {
                parentID = id;
            }
        }

        this.mainID = mainID;
        this.parentID = parentID;
    }

   /**
    * Casts Component Matcher to an instance of type SeleComponentID
    * @param componentID the component ID
    * @return instance of {@link com.vmware.suitaf.apl.sele.SeleComponentID}
    */
   public static SeleComponentID fromMatcher(ComponentMatcher componentID) {
      Preconditions.checkNotNull(componentID);
      return (SeleComponentID) componentID;
   }

   /**
    * Return the Direct ID of the component matcher
    * @param componentMatcher the component matcher
    * @return the direct ID
    */
   public static DirectID toDirectId(ComponentMatcher componentMatcher) {
      return fromMatcher(componentMatcher).mainID;
   }

   @Override
    public String getLogForm() {
        return toString();
    }

    @Override
    public String toString() {
        return
        ((mainID == null)? "": (mainID + ";")) +
        ((parentID == null)? "": (parentID + ";"));
    }
}
