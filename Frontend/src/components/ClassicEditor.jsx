// Import ImageResize plugin
import ImageResize from "@ckeditor/ckeditor5-image/src/imageresize";

ClassicEditor.create(document.querySelector("#editor"), {
  plugins: [Image, ImageResize /* ... other plugins ... */],
  image: {
    // Konfigurasi untuk ImageResize
    resizeOptions: [
      {
        name: "resizeImage:original",
        value: null,
        icon: "original",
      },
      {
        name: "resizeImage:50",
        value: "50",
        icon: "medium",
      },
      {
        name: "resizeImage:75",
        value: "75",
        icon: "large",
      },
    ],
    toolbar: [
      "imageTextAlternative",
      "|",
      "resizeImage:50",
      "resizeImage:75",
      "resizeImage:original",
    ],
  },
  // ... konfigurasi lainnya ...
}).catch((error) => {
  console.error(error);
});
