package com.vmware.hsua.common.datamodel;

import java.lang.annotation.Documented;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;


/**
 * This class represents DEMO-specific extension of the base data container
 * class {@link PropertyBox}. It adds annotations that allow attaching
 * meta-data to each property field.<br>
 * The class also adds functionality for better property-value setting and
 * easier usage.
 *
 * @author dkozhuharov
 */
public class BasePBox extends PropertyBox {

    /**
     * Enumerator describing the set of supported languages for I18N and I10N.
     */
    public static enum Lang {
        JP("\u30E6\u30B6\u3078\u305D"),
        DE("\u00FC\u00DF\u00E4\u00F6"),
        IT("\u00E0\u00E8\u00E9\u00EC"),
        FR("\u00F9\u00FB\u00FC\u00FF"),
        EN("");

        /**
         * Holds UTF encoded symbols that are specific for that language.
         * The symbols are used for I18N testing.
         */
        public final String i18nSuffix;

        private Lang(String i18nString) {
            this.i18nSuffix = i18nString;
        }

        private static BasePBox.Lang current = EN;
        /**
         * Sets the current language for the whole UI automation projects.
         * This controls the behavior of all I18N and I10N features
         * @param newLang - the new current language
         */
        public static void setCurrent(BasePBox.Lang newLang) {
            current = newLang;
        }

        /**
         * Return the language for which this test case is run.
         * @return
         */
        public static BasePBox.Lang getCurrent() {
            return current;
        }
    }

    /**
     * This property container is DEMO-application specific extension of the
     * basic abstract property. It adds property-level meta data annotations.
     * Also this property container has additional functionality relevant to
     * the annotation's setting and getting actions.
     * @author dkozhuharov
     *
     * @param <E> - the type of the value to be contained in the property
     */
    public static class DataProperty<E> extends AbstractProperty<E> {
        @SuppressWarnings("unchecked")
        @Override
        protected E setTransformer(E value) {
            // Implement I18N related transformation at "set" time
            if (this.annotation(BasePBox.I18N.class) != null
                    && value instanceof String) {
                return (E) (((String) value) + Lang.getCurrent().i18nSuffix);
            }
            return value;
        }
    }

    /**
     * This method returns the String value of the property field annotated
     * with the {@link Name} annotation.
     *
     * @return the unique name of the entity represented by this PBox
     */
    public final String name() {
        AbstractProperty<?> name = getAnnotated(Name.class);
        if (name != null) {
         return name.get().toString();
      } else {
         return null;
      }
    }

    /**
     * This annotation is to be used in the {@link BasePBox} classes to
     * mark the I18N compliant fields.
     * <br><br>
     * When in I18N test mode - the values of all the compliant properties will
     * be expanded automatically with I18N specific symbols.
     *
     * @author dkozhuharov
     *
     */
    @Retention(RetentionPolicy.RUNTIME)
    @Documented
    public static @interface I18N {}

    /**
     * This annotation is to be used in the {@link BasePBox} classes to
     * mark the field that carries the unique name of the entity.
     * <br><br>
     * The marked fields will be automatically used in the methods that check
     * presence of the entity on screen and check its status in grid.
     *
     * @author dkozhuharov
     */
    @Retention(RetentionPolicy.RUNTIME)
    @Documented
    public static @interface Name {}
}