# Imaginarium

A simple golang image storage engine. Used to create and store different sizes/thumbnails of user uploaded images.

## Description

**Imaginarium** enables you to create copies (or thumbnails) of your images and stores
them along with the original image on your filesystem. The image and its
copies are stored in a file structure based on the MD5 checksum of the original image.
The first character of the checksum used as the lvl 1 directory name.

**Imaginarium** supports png, jpg formats. The decoder for any given image is chosen by the image's mimetype.

## Usage
```html
<form action="http://localhost/upload" method="post" enctype="multipart/form-data">
    Context: <label><input type="text" name="contexts"></label>
    Files: <input type="file" name="image" multiple>
    <input type="submit" value="Submit">
</form>
