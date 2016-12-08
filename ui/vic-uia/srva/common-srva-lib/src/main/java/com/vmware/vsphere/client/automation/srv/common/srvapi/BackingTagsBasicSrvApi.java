/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.net.URI;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.sso.SsoClient;
import com.vmware.client.automation.util.SsoUtil;
import com.vmware.vim.binding.dataservice.QName;
import com.vmware.vim.binding.dataservice.tagging.Category;
import com.vmware.vim.binding.dataservice.tagging.CategoryInfo;
import com.vmware.vim.binding.dataservice.tagging.CategoryInfo.Cardinality;
import com.vmware.vim.binding.dataservice.tagging.Tag;
import com.vmware.vim.binding.dataservice.tagging.TagInfo;
import com.vmware.vim.binding.dataservice.tagging.TagManager;
import com.vmware.vim.binding.impl.dataservice.QNameImpl;
import com.vmware.vim.binding.impl.dataservice.tagging.CategoryInfoImpl;
import com.vmware.vim.binding.impl.dataservice.tagging.TagInfoImpl;
import com.vmware.vim.binding.vmodl.ManagedObject;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vim.query.client.Client;
import com.vmware.vim.vmomi.cis.CisIdConverter;
import com.vmware.vim.vmomi.core.impl.BlockingFuture;
import com.vmware.vise.search.impl.NamespaceUtils;
import com.vmware.vise.util.xml.NamespaceUtil;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingCategorySpec;
import com.vmware.vsphere.client.automation.srv.common.spec.BackingTagSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.UpdateBackingTagSpec;

/**
 * API commands for CRUD operation on backing tags and categories.
 */
public class BackingTagsBasicSrvApi {

   private static final String NAMESPACE = "vim25";
   private static final String DEFAULT_NAMESPACE_PREFIX = "urn:";
   private static final Logger _logger = LoggerFactory
         .getLogger(BackingTagsBasicSrvApi.class);

   private static BackingTagsBasicSrvApi instance = null;

   protected BackingTagsBasicSrvApi() {
   }

   /**
    * Get instance of BackingTagsSrvApi.
    *
    * @return created instance
    */
   public static BackingTagsBasicSrvApi getInstance() {
      if (instance == null) {
         synchronized (BackingTagsBasicSrvApi.class) {
            if (instance == null) {
               _logger.info("Initializing BackingTagsSrvApi.");
               instance = new BackingTagsBasicSrvApi();
            }
         }
      }

      return instance;
   }

   /**
    * Creates a backing category.
    *
    * @param category
    *           - the spec for the category
    * @return the ID of the newly created tag category
    * @throws Exception
    *            if can not get access to the TagManager API
    */
   public boolean createBackingCategory(BackingCategorySpec category)
         throws Exception {
      validateSpec(category);

      com.vmware.vim.vmomi.core.Future<ManagedObjectReference> f = new BlockingFuture<ManagedObjectReference>();
      CategoryInfo categoryInfo = categorySpecToInfo(category);

      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         queryClient.getTagManager().createCategory(categoryInfo, f);
         return f.get() != null;
      } finally {
         queryClient.close();
      }
   }

   /**
    * Deletes a backing category safely
    *
    * @param category
    *           - the spec of the category to delete
    * @throws Exception
    *            if can not get access to the TagManager API
    */
   public boolean deleteBackingCategorySafely(BackingCategorySpec category)
         throws Exception {
      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         ManagedObjectReference categoryMoRef = getBackingCategory(queryClient,
               category);
         if (categoryMoRef != null) {
            Category c = createStub(queryClient, Category.class, categoryMoRef);
            com.vmware.vim.vmomi.core.Future<Void> f = new BlockingFuture<Void>();
            c.delete(f);
            try {
               f.get();
               return true;
            } catch (Exception e) {
               e.printStackTrace();
               return false;
            }
         }

         return false;
      } finally {
         queryClient.close();
      }
   }

   /**
    * Checks if a backing category is existing.
    *
    * @param category
    *           - the backing category
    * @return true if the backing category exists
    * @throws Exception
    *            if can not get access to the API service
    */
   public boolean checkBackingCategoryExists(BackingCategorySpec category)
         throws Exception {
      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         return getBackingCategory(queryClient, category) != null;
      } finally {
         queryClient.close();
      }
   }

   /**
    * Returns a list with all backing categories
    *
    * @param listNotUsed
    *           - if we want to list the unused categories
    * @param listUsed
    *           - if we want to list the used categories
    * @return a list with all backing categories
    * @throws Exceptionif
    *            can not get access to the API service
    */
   public List<BackingCategorySpec> listBackingCategories(
         BackingCategorySpec category, boolean listNotUsed, boolean listUsed)
         throws Exception {
      List<BackingCategorySpec> result = new ArrayList<BackingCategorySpec>();
      // TODO: fix me to work with service spec
      // SsoClient ssoConnection = SsoUtil.getConnector(new ServiceSpec())
      // .getConnection();
      // Client queryClient = ssoConnection.getQueryClient();
      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         TagManager tagManager = queryClient.getTagManager();

         ManagedObjectReference[] categories = tagManager.enumerateCategories();
         if (categories != null) {
            for (ManagedObjectReference categoryMoRef : categories) {
               Category c = createStub(queryClient, Category.class,
                     categoryMoRef);
               CategoryInfo catInfo = c.getInfo();
               if ((listNotUsed && (catInfo.getUsedBy() == null))
                     || (listUsed && (catInfo.getUsedBy() != null))) {
                  result.add(categoryInfoToSpec(catInfo));
               }
            }
         }
         return result;
      } finally {
         queryClient.close();
      }
   }

   /**
    * Creates a backing tag
    *
    * @param category
    *           - the backing category that will contain the tag
    * @param tag
    *           - the spec of the tag
    * @return the ID of the newly created tag
    * @throws Exception
    *            if we can't get access to the service
    */
   public boolean createBackingTag(BackingCategorySpec category,
         BackingTagSpec tag) throws Exception {
      validateSpec(tag);
      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         ManagedObjectReference categoryMor = getBackingCategory(queryClient,
               category);
         TagInfo tagInfo = tagSpecToInfo(tag);
         Category cat = createStub(queryClient, Category.class, categoryMor);
         com.vmware.vim.vmomi.core.Future<ManagedObjectReference> f = new BlockingFuture<ManagedObjectReference>();
         cat.createTag(tagInfo, f);
         return f.get() != null;
      } finally {
         queryClient.close();
      }
   }

   /**
    * Update backing tag that belongs to a specific category.
    *
    * @param updateSpec
    *           spec that contains all necessary info for updating the tag
    * @return true if the operation is successful, false otherwise
    * @throws Exception
    *            if we can't get access to the service
    */
   public boolean updateBackingTag(UpdateBackingTagSpec updateSpec)
         throws Exception {
      BackingCategorySpec category = updateSpec.category.get();
      BackingTagSpec targetTag = updateSpec.targetTag.get();
      BackingTagSpec newTag = updateSpec.newTargetConfigs.get();

      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         ManagedObjectReference categoryMor = getBackingCategory(queryClient,
               category);
         if (categoryMor != null) {
            ManagedObjectReference tagMoRef = getBackingTag(queryClient,
                  categoryMor, targetTag);
            Tag tag = createStub(queryClient, Tag.class, tagMoRef);
            // invoking update backing tag asynchronously
            if (tagMoRef != null) {
               com.vmware.vim.vmomi.core.Future<Void> f = new BlockingFuture<Void>();
               tag.updateInfo(tagSpecToInfo(newTag), f);
               try {
                  f.get();
                  return true;
               } catch (Exception e) {
                  _logger.error(e.getMessage());
                  _logger.error(Arrays.asList(e.getStackTrace()).toString());
                  return false;
               }
            } else {
               _logger.error(String.format("Tag c%s was not found!",
                     targetTag.name.get()));
            }
         } else {
            _logger.error(String.format("Tag category %s was not found!",
                  category.name.get()));
         }
         return false;
      } finally {
         queryClient.close();
      }
   }

   /**
    * Deletes a backing tag.
    *
    * @param category
    *           - the backing category that will contain the tag
    * @param tag
    *           - the spec of the tag to delete
    * @throws Exception
    *            - if we can't access the API
    */
   public boolean deleteBackingTagSafely(BackingCategorySpec category,
         BackingTagSpec tag) throws Exception {
      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         ManagedObjectReference categoryMor = getBackingCategory(queryClient,
               category);
         if (categoryMor != null) {
            ManagedObjectReference tagMoRef = getBackingTag(queryClient,
                  categoryMor, tag);
            if (tagMoRef != null) {
               // Delete
               Tag t = createStub(queryClient, Tag.class, tagMoRef);
               com.vmware.vim.vmomi.core.Future<Void> f = new BlockingFuture<Void>();
               t.delete(f);
               try {
                  f.get();
                  return true;
               } catch (Exception e) {
                  e.printStackTrace();
                  return false;
               }
            }
         }

         return false;
      } finally {
         queryClient.close();
      }
   }

   /**
    * Checks if a backing tag is existing.
    *
    * @param category
    *           - the backing category that will contain the tag
    * @param tag
    *           - the backing tag
    * @return true if the backing tag exists
    * @throws Exception
    *            if can not get access to the API service
    */
   public boolean checkBackingTagExists(BackingCategorySpec category,
         BackingTagSpec tag) throws Exception {
      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         ManagedObjectReference categoryMor = getBackingCategory(queryClient,
               category);
         if (categoryMor != null) {
            return getBackingTag(queryClient, categoryMor, tag) != null;
         }
         return false;
      } finally {
         queryClient.close();
      }
   }

   /**
    * Attaches resources to the given tag.
    *
    * @param category
    *           the category the tag belongs to
    * @param tag
    *           the tag to attach resources to
    * @param resourceSpecs
    *           a list of entities to assign the tag to
    * @throws Exception
    *            when an API call returns exception
    */
   public void attachResources(BackingCategorySpec category,
         BackingTagSpec tag, List<? extends ManagedEntitySpec> resourceSpecs)
         throws Exception {

      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();

      try {
         // Find the tag
         ManagedObjectReference categoryMoRef = getBackingCategory(queryClient,
               category);
         ManagedObjectReference tagMoRef = getBackingTag(queryClient,
               categoryMoRef, tag);

         // Get the URIs of the resources
         List<URI> resourceUris = ManagedEntityUtil
               .getResourceUris(resourceSpecs);

         // Make the "attach" call
         TagManager tagManager = queryClient.getTagManager();
         for (URI resourceUri : resourceUris) {
            tagManager.attachTagsToObject(resourceUri,
                  new ManagedObjectReference[] { tagMoRef });
         }
      } finally {
         queryClient.close();
      }
   }

   /**
    * Detaches resources from the given tag.
    *
    * @param category
    *           the category the tag belongs to
    * @param tag
    *           the tag to detach resources from
    * @param resourceSpecs
    *           a list of resources to attach
    * @throws Exception
    *            when an API call returns exception
    */
   public void detachResources(BackingCategorySpec category,
         BackingTagSpec tag, List<? extends ManagedEntitySpec> resourceSpecs)
         throws Exception {

      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();
      try {
         // Find the tag
         ManagedObjectReference categoryMoRef = getBackingCategory(queryClient,
               category);
         ManagedObjectReference tagMoRef = getBackingTag(queryClient,
               categoryMoRef, tag);

         // Get the URIs of the resources
         List<URI> resourceUris = ManagedEntityUtil
               .getResourceUris(resourceSpecs);

         // Make the "detach" call
         TagManager tagManager = queryClient.getTagManager();
         for (URI resourceUri : resourceUris) {
            tagManager.detachTagsFromObject(resourceUri,
                  new ManagedObjectReference[] { tagMoRef });
         }
      } finally {
         queryClient.close();
      }
   }

   /**
    * Lists attached resources to given tag.
    *
    * @param category
    *           the category the tag belongs to
    * @param tag
    *           the tag resources are attached to
    * @return a list of MoRefs of resources assigned to the tag
    * @throws Exception
    *            when an API call returns exception
    */
   public List<ManagedObjectReference> listAttachedResources(
         BackingCategorySpec category, BackingTagSpec tag) throws Exception {

      List<ManagedObjectReference> resourceMoRefs = new ArrayList<ManagedObjectReference>();

      SsoClient ssoConnection = SsoUtil.getVcConnector(category)
            .getConnection();
      Client queryClient = ssoConnection.getQueryClient();

      try {
         // Find the tag
         ManagedObjectReference categoryMoRef = getBackingCategory(queryClient,
               category);
         ManagedObjectReference tagMoRef = getBackingTag(queryClient,
               categoryMoRef, tag);

         // Make the "query attached resources" call
         TagManager tagManager = queryClient.getTagManager();
         URI[] objUris = tagManager.queryAttachedObjects(tagMoRef);

         if (objUris != null) {
            for (URI attachedObject : objUris) {
               String uri = attachedObject.toASCIIString();
               ManagedObjectReference moRef = CisIdConverter.fromCisId(uri);
               resourceMoRefs.add(moRef);
            }
         }
      } finally {
         queryClient.close();
      }

      return resourceMoRefs;
   }

   /**
    * Check if tag is assigned to resource.
    *
    * @param category
    *           the category the tag belongs to
    * @param tag
    *           the tag which resources are being examined
    * @param resource
    *           the resource that is being looked for
    * @return true if tag is attached to resource, otherwise - false
    * @throws Exception
    *            when an API call returns exception
    */
   public boolean checkTagAttachedToResource(BackingCategorySpec category,
         BackingTagSpec tag, ManagedEntitySpec resource) throws Exception {

      // List all resources attached to given tag
      List<ManagedObjectReference> resourceMoRefs = listAttachedResources(
            category, tag);

      // Search for given resource by name
      for (ManagedObjectReference resourceMoRef : resourceMoRefs) {
         // It is possible a resource to be deleted before getting its info.
         // In such case a ManagedObjectNotFound will be thrown.
         try {
            if (ManagedEntityUtil.getNameFromMoRef(resourceMoRef,
                  resource.service.get()).equals(resource.name.get())) {
               return true;
            }
         } catch (com.vmware.vim.binding.vmodl.fault.ManagedObjectNotFound monfe) {
            _logger.error(monfe.getMessage());
         }

      }

      // Resource was not found in list
      return false;
   }

   // ---------------------------------------------------------------------------
   // Package methods

   public ManagedObjectReference getBackingCategory(Client queryClient,
         ManagedEntitySpec category) throws Exception {
      TagManager tagManager = queryClient.getTagManager();
      ManagedObjectReference[] categories = tagManager.enumerateCategories();
      if (categories != null) {
         for (ManagedObjectReference categoryMoRef : categories) {
            Category c = createStub(queryClient, Category.class, categoryMoRef);
            // It is possible a category to be deleted before getting its info.
            // In such case a ManagedObjectNotFound will be thrown.
            try {
               CategoryInfo catInfo = c.getInfo();
               if (catInfo.getName().equals(category.name.get())) {
                  return categoryMoRef;
               }
            } catch (com.vmware.vim.binding.vmodl.fault.ManagedObjectNotFound monfe) {
               _logger.error(monfe.getMessage());
            }
         }
      }
      return null;
   }

   public ManagedObjectReference getBackingTag(Client queryClient,
         ManagedObjectReference categoryMor, BackingTagSpec tag)
         throws Exception {
      TagManager tagManager = queryClient.getTagManager();

      ManagedObjectReference[] tags = tagManager
            .queryTagsByCategory(categoryMor);
      if (tags != null) {
         for (ManagedObjectReference tagMoRef : tags) {
            Tag t = createStub(queryClient, Tag.class, tagMoRef);
            // It is possible a tag to be deleted before getting its info.
            // In such case a ManagedObjectNotFound will be thrown.
            try {
               TagInfo tagInfo = t.getInfo();
               if (tagInfo.getName().equals(tag.name.get())) {
                  return tagMoRef;
               }
            } catch (com.vmware.vim.binding.vmodl.fault.ManagedObjectNotFound monfe) {
               _logger.error(monfe.getMessage());
            }
         }
      }
      return null;
   }

   // ---------------------------------------------------------------------------
   // Private methods

   private void validateSpec(ManagedEntitySpec spec) {
      if (spec == null) {
         throw new IllegalArgumentException("Spec is not set");
      }

      if (!spec.name.isAssigned()) {
         throw new IllegalArgumentException("Name is not set");
      }
   }

   /**
    * Create a stub of type clazz for the passed-in moRef.
    *
    * @throws Exception
    */
   private <T extends ManagedObject> T createStub(Client queryClient,
         Class<T> clazz, ManagedObjectReference moRef) throws Exception {

      T t = queryClient.getVmomiClient().createStub(clazz, moRef);
      return t;
   }

   private CategoryInfo categorySpecToInfo(BackingCategorySpec category) {
      String desc = category.description.isAssigned() ? category.description
            .get() : null;
      Cardinality cardinality = !category.cardinality.isAssigned() ? CategoryInfo.Cardinality.single
            : category.cardinality.get();
      List<String> associableObjects = category.associatedObjects.getAll();
      List<QName> associableObjectsInfo = new LinkedList<>();
      String namespace;
      String unqualifiedName;
      for (String associableObject : associableObjects) {
         namespace = DEFAULT_NAMESPACE_PREFIX + NAMESPACE;
         unqualifiedName = NamespaceUtil
               .getUnqualifiedTypeName(associableObject);
         associableObjectsInfo.add(new QNameImpl(namespace, unqualifiedName));
      }

      QName[] array = associableObjectsInfo
            .toArray(new QName[associableObjectsInfo.size()]);

      CategoryInfo categoryInfo = new CategoryInfoImpl(category.name.get(),
            desc, cardinality, array.length == 0 ? null : array,
            null);
      return categoryInfo;
   }

   private BackingCategorySpec categoryInfoToSpec(CategoryInfo catInfo) {
      BackingCategorySpec result = new BackingCategorySpec();
      result.name.set(catInfo.getName());
      result.description.set(catInfo.getDescription());
      return result;
   }

   private TagInfo tagSpecToInfo(BackingTagSpec tag) {
      String desc = tag.description.isAssigned() ? tag.description.get() : null;
      TagInfo tagInfo = new TagInfoImpl(tag.name.get(), desc, new String[] {});
      return tagInfo;
   }

}
