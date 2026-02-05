class UploadAdapter {
  constructor(loader) {
    // Simpan referensi ke image loader
    this.loader = loader;
  }

  // Mulai mengunggah gambar
  upload() {
    return this.loader.file.then(
      (file) =>
        new Promise((resolve, reject) => {
          this._initRequest();
          this._initListeners(resolve, reject);
          this._sendRequest(file);
        })
    );
  }

  // Inisialisasi request XHR
  _initRequest() {
    const xhr = (this.xhr = new XMLHttpRequest());
    xhr.open(
      "POST",
      "http://localhost:8080/api-blog-ngebruk/upload_image",
      true
    );
    xhr.responseType = "json";
  }

  // Inisialisasi listeners untuk XHR request
  _initListeners(resolve, reject) {
    const xhr = this.xhr;
    const loader = this.loader;
    const genericErrorText = `Couldn't upload file: ${loader.file.name}.`;

    xhr.addEventListener("error", () => reject(genericErrorText));
    xhr.addEventListener("abort", () => reject());
    xhr.addEventListener("load", () => {
      const response = xhr.response;
      if (!response || !response.presigned_urls) {
        return reject(
          response && response.error ? response.error.message : genericErrorText
        );
      }
      resolve({ default: response.presigned_urls });
    });

    // Support for token atau header lain jika diperlukan
    // if (this.token) {
    //   xhr.setRequestHeader("Authorization", "Bearer " + this.token);
    // }
  }

  // Mengirim request
  _sendRequest(file) {
    const data = new FormData();
    data.append("images", file);
    this.xhr.send(data);
  }
}

// Helper untuk menambahkan adapter ke CKEditor
function UploadAdapterPlugin(editor) {
  editor.plugins.get("FileRepository").createUploadAdapter = (loader) => {
    return new UploadAdapter(loader);
  };
}
export default UploadAdapterPlugin;
