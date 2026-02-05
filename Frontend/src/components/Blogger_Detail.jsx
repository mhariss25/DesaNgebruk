import { useState, useEffect } from "react";
import axios from "axios";
import { useParams, Link } from "react-router-dom";
import { format } from "date-fns";
import DOMPurify from "dompurify";

const Blogger_Detail = () => {
  const { id } = useParams();
  const [hasil, setHasil] = useState(null);
  const [recentBlogs, setRecentBlogs] = useState([]);

  // Function to fetch main blog data
  const fetchData = async () => {
    try {
      const response = await axios.get(
        `http://localhost:8080/api-blog-ngebruk/blogger/${id}`
      );
      setHasil(response.data);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  // Function to fetch recent blogs data
  const fetchRecentBlogs = async () => {
    try {
      const response = await axios.get(
        `http://localhost:8080/api-blog-ngebruk/blogger?pageSize=5`
      );
      // Sort the blogs by created_at in descending order
      const sortedRecentBlogs = response.data.bloggers.sort(
        (a, b) => new Date(b.created_at) - new Date(a.created_at)
      );
      setRecentBlogs(sortedRecentBlogs);
    } catch (error) {
      console.error("Error fetching recent blogs:", error);
    }
  };

  useEffect(() => {
    fetchData();
  }, [id]);

  useEffect(() => {
    fetchRecentBlogs();
  }, []);

  const formatTanggal = (created_at) => {
    return format(new Date(created_at), " dd MMMM yyyy HH:mm:ss");
  };

  const createMarkup = (htmlContent) => {
    return {
      __html: DOMPurify.sanitize(htmlContent),
    };
  };

  const stripHtml = (html) => {
    const tempDiv = document.createElement("div");
    tempDiv.innerHTML = html;
    return tempDiv.textContent || tempDiv.innerText || "";
  };

  const sliceword = (description) => {
    const words = description.split(" ");
    const first8Words = words.slice(0, 10).join(" ");

    return first8Words;
  };

  return (
    <div className="p-8 mt-16 mb-40 flex-grow ">
      <div className="flex flex-col md:flex-row">
        <div className="flex-grow max-w-6xl">
          {hasil && (
            <>
              <div className="text-2xl font-bold text-emerald-500">
                {hasil.name_blog}
              </div>
              <img
                className="w-full h-56 object-cover"
                src={hasil.heading_bloger}
                alt=""
              />
              <div className="mt-4">
                <p className="text-lg font-bold text-gray-500">
                  {hasil.user.nama} |
                  <span className=" text-base font-semibold">
                    {formatTanggal(hasil.created_at)}
                  </span>{" "}
                  |{" "}
                  <span className="text-base font-semibold">
                    {hasil.kategori.kategori_name}
                  </span>
                </p>
                <div
                  className="font-serif text-gray-700 prose prose-img:rounded-xl prose-headings:underline prose-w-full prose-a:text-blue-600"
                  dangerouslySetInnerHTML={createMarkup(hasil.fill_blogger)}
                ></div>
              </div>
            </>
          )}
        </div>

        <div className="w-full md:w-2/5 max-w-md md:ml-4 mt-4 md:mt-2 bg-gray-50 shadow-lg shadow-slate-300 border-2 rounded-lg p-3">
          <h2 className="text-xl font-bold mb-4">Konten Terbaru</h2>
          <ul>
            {recentBlogs.map((blog) => (
              <li
                key={blog.id_blogger}
                className="mb-3 p-3 hover:bg-slate-200 hover:bg-opacity-40 border rounded-lg shadow-md shadow-slate-300 border-gray-300"
              >
                <Link to={`/blogger_detail/${blog.id_blogger}`}>
                  <h3 className="text-lg font-semibold">{blog.name_blog}</h3>
                  <p className=" text-sm">
                    {stripHtml(sliceword(blog.fill_blogger))}
                  </p>
                  <p className="text-sm text-gray-600">
                    {formatTanggal(blog.created_at)}
                  </p>
                </Link>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
};
export default Blogger_Detail;
