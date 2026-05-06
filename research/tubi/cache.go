package main

import (
   "encoding/xml"
   "fmt"
   "log"
   "os"
   "path/filepath"
   "reflect"
)

// Cache holds the pre-computed OS path for the cache directory.
type Cache struct {
   FullPath string
}

// Setup computes the full cache path, creates the directory exactly once,
// and stores the path in the Cache struct.
func (c *Cache) Setup(dirName string) error {
   cacheDir, err := os.UserCacheDir()
   if err != nil {
      return fmt.Errorf("failed to get cache directory: %w", err)
   }

   c.FullPath = filepath.Join(cacheDir, dirName)

   // Create the directory immediately upon setup
   if err := os.MkdirAll(c.FullPath, os.ModePerm); err != nil {
      return fmt.Errorf("failed to create directory: %w", err)
   }

   return nil
}

// GetFilePath unwraps pointers and builds the absolute string path for the file.
// Exported so users can manually locate, check, or delete cache files.
func (c *Cache) GetFilePath(v any) string {
   t := reflect.TypeOf(v)
   for t.Kind() == reflect.Ptr {
      t = t.Elem()
   }

   return filepath.Join(c.FullPath, t.Name()+".xml")
}

// Encode marshals the value and writes it to the cache directory.
// It logs the full path immediately before attempting to write the file.
func (c *Cache) Encode(v any) error {
   filename := c.GetFilePath(v)

   data, err := xml.MarshalIndent(v, "", "  ")
   if err != nil {
      return fmt.Errorf("failed to encode XML: %w", err)
   }

   log.Printf("Creating file: %s\n", filename)

   err = os.WriteFile(filename, data, os.ModePerm)
   if err != nil {
      return fmt.Errorf("failed to write file: %w", err)
   }

   return nil
}

// Decode reads the XML from the cache directory and populates the struct.
func (c *Cache) Decode(v any) error {
   filename := c.GetFilePath(v)

   data, err := os.ReadFile(filename)
   if err != nil {
      return fmt.Errorf("failed to read file: %w", err)
   }

   err = xml.Unmarshal(data, v)
   if err != nil {
      return fmt.Errorf("failed to decode XML: %w", err)
   }

   return nil
}
