import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Cookies from "js-cookie";
import { CKEditor } from "@ckeditor/ckeditor5-react";
import ClassicEditor from "@ckeditor/ckeditor5-build-classic";
import UploadAdapterPlugin from "../UploadAdapter";
import { useParams } from "react-router-dom";
import Swal from "sweetalert2";

const FormBloggerAdmin = () => {
  const { id_blogger } = useParams();
  const navigate = useNavigate();
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState("");
  const [formData, setFormData] = useState({
    name_blog: "",
    fill_blogger: "",
    heading_blogger: null,
  });

  useEffect(() => {
    axios
      .get("http://localhost:8080/api-blog-ngebruk/kategori")
      .then((response) => {
        setCategories(response.data);
      })
      .catch((error) => {
        console.error("Error fetching categories:", error);
      });

    if (id_blogger) {
      axios
        .get(`http://localhost:8080/api-blog-ngebruk/blogger/${id_blogger}`)
        .then((response) => {
          setFormData({ ...response.data });
          setSelectedCategory(response.data.kategori_id.toString());
          window.dispatchEvent(new Event("userUpdated"));
        })
        .catch((error) => {
          console.error("Error fetching blogger data:", error);
        });
    }
  }, [id_blogger]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;

    if (name === "name_blog" && value.length > 99) {
      return;
    }
    setFormData((prevState) => ({
      ...prevState,
      [name]: value,
    }));
  };

  const handleEditorChange = (event, editor) => {
    const data = editor.getData();
    setFormData((prevState) => ({ ...prevState, fill_blogger: data }));
  };

  const handleFileChange = (e) => {
    setFormData((prevState) => ({
      ...prevState,
      heading_blogger: e.target.files[0],
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    const token = Cookies.get("token");
    const config = {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "multipart/form-data",
      },
    };

    const data = new FormData();
    data.append("name_blog", formData.name_blog);
    data.append("fill_blogger", formData.fill_blogger);
    data.append("kategori_id", selectedCategory);
    if (formData.heading_blogger) {
      data.append("heading_blogger", formData.heading_blogger);
    }

    const url = id_blogger
      ? `http://localhost:8080/api-blog-ngebruk/blogger/${id_blogger}`
      : `http://localhost:8080/api-blog-ngebruk/CreateBlogger`;

    const method = id_blogger ? axios.patch : axios.post;

    method(url, data, config)
      .then((response) => {
        console.log("Blogger operation successful:", response.data);
        // Tampilkan Sweet Alert
        Swal.fire({
          title: id_blogger ? "Updated Successfully!" : "Added Successfully!",
          text: "Your blogger data has been saved.",
          icon: "success",
        }).then(() => {
          navigate("/blogger/listdashboardadmin");
        });
      })
      .catch((error) => {
        console.error("Error in blogger operation:", error);
        console.error("Error in blogger operation:", error.response.data);
        // Tampilkan Sweet Alert untuk error
        swal.fire("Error", "There was an error saving your data", "error");
      });
  };

  const handleCategoryChange = (e) => {
    setSelectedCategory(e.target.value);
  };

  return (
    <>
      <div className=" p-4 border-2 rounded-lg  mt-20">
        <section className="bg-white dark:bg-gray-900">
          <div className="py-8 px-4 mx-auto max-w-screen-md lg:py-16">
            <h2 className="mb-4 text-xl text-center sm:text-left font-bold text-gray-900 dark:text-white">
              {id_blogger ? "Edit Konten" : "Add Konten"}
            </h2>
            <form encType="multipart/form-data" onSubmit={handleSubmit}>
              <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 sm:gap-6">
                <div className="sm:col-span-2">
                  <label
                    htmlFor="name"
                    className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                  >
                    Judul Konten
                  </label>
                  <input
                    type="text"
                    name="name_blog"
                    value={formData.name_blog}
                    onChange={handleInputChange}
                    id="name"
                    className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                    placeholder="Masukkan Judul Konten"
                    required=""
                  />
                  <p className="text-sm text-gray-500">
                    {formData.name_blog.length}/99 characters
                  </p>
                </div>
                <div className="w-full">
                  <label
                    htmlFor="brand"
                    className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                  >
                    Kategori
                  </label>
                  <select
                    className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                    value={selectedCategory}
                    onChange={handleCategoryChange}
                  >
                    <option value="">Pilih Kategori</option>
                    {categories.map((category) => (
                      <option
                        key={category.id_kategori}
                        value={category.id_kategori}
                      >
                        {category.kategori_name}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="w-full">
                  <label className="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
                    Cover Images
                  </label>
                  <input
                    className="block w-full text-sm text-gray-900 border border-gray-300 rounded-lg cursor-pointer bg-gray-50 dark:text-gray-400 focus:outline-none dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400"
                    name="heading_blogger"
                    onChange={handleFileChange}
                    id="user_avatar"
                    type="file"
                  />{" "}
                  <img src={formData.heading_bloger} />
                </div>

                <div className="sm:col-span-2 prose-ul:list-disc">
                  <label
                    htmlFor="description"
                    className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
                  >
                    Isi Konten
                  </label>
                  <CKEditor
                    value={formData.fill_blogger}
                    editor={ClassicEditor}
                    data={formData.fill_blogger}
                    onChange={handleEditorChange}
                    config={{
                      extraPlugins: [UploadAdapterPlugin],
                      // Konfigurasi lainnya
                    }}
                  />
                </div>
              </div>
              <button
                type="submit"
                className="inline-flex mb-32 items-center px-5 py-2.5 mt-4 sm:mt-6 text-sm font-medium text-center text-white bg-blue-700 rounded-lg focus:ring-4 focus:ring-blue-200 dark:focus:ring-blue-900 hover:bg-blue-800"
              >
                {id_blogger ? "Edit Konten" : "Tambah Konten"}
              </button>
            </form>
          </div>
        </section>
      </div>
    </>
  );
};

export default FormBloggerAdmin;
