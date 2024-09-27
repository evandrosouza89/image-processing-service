<p align="center">
    <img src="/assets/screenshot.png" alt="screenshot">
</p>

# Image processing service

Image processing service is a demo web application that allows users to upload multiple images and perform batch operations like rotating, resizing, and enlarging on them. The backend of this project is built using Go, and the frontend is designed using HTML, CSS, and JavaScript.

## Table of Contents
- [Installation](#features)
- [Usage](#usage)
- [Technology Stack](#stack)
- [API Reference](#api)
- [Explanation of Key UI Elements](#key)
- [License](#license)

<a id="features"></a>
## Features
- **Drag and Drop**: Users can drag and drop images into the dropzone.
- **Batch Operations**: Choose from multiple operations (rotate, resize, enlarge) to apply to all selected images.
- **Preview**: Preview the uploaded images before performing the operation.
- **Download**: Download the processed images in batch.

<a id="usage"></a>
## Usage
1. Open the application in your browser.
2. Drag and drop images into the dropzone (Maximum 3 images, 350KB each).
3. Choose an operation from the dropdown:
    - **Rotate**: Rotate the images by 90, 180, or 270 degrees.
    - **Resize**: Set the desired width and height to resize the images.
    - **Enlarge**: Enlarge the images by a percentage preserving its ratio. Must be between 1% and 200%, excluding 100%.
4. Click the Submit button to process the images.
5. Download the processed images.

<a id="stack"></a>
## Technology Stack
- **Frontend**: Bootstrap 
- **Backend**: Go
- **Image Processing**: Bimg
- **Containerization**: Docker

<a id="api"></a>
## API Reference
### Endpoints
- ```POST /upload```: Accepts images and operation details, processes the images, and returns the modified images.
### Request
- ```Content-Type```: multipart/form-data
- **Parameters**:
  - ```images``` (array of images): The images to be processed.
  - ```operation``` (JSON object): The operation and its parameters, e.g., ```{ "operation": "rotate", "parameters": { "angle": 90 } }```.
### Response
 - **Content-Type**: ```multipart/form-data```
 - The processed images are returned in the response.

<a id="key"></a>
### Explanation of Key UI Elements
- **Dropzone**: Users can drag and drop images here.
- **Preview Section**: Displays a preview of the uploaded images.
- **Operation Form**: Allows the user to choose and configure operations.
- **Submit Button**: Sends the images and operation parameters to the backend for processing.

<a id="license"></a>
## License
**image-processing-service** is licensed under the MIT License. A short and simple permissive license with conditions only requiring preservation of copyright and license notices. Licensed works, modifications, and larger works may be distributed under different terms and without source code.
