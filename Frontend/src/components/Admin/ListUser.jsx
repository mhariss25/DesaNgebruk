import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Cookies from "js-cookie";
import { format } from "date-fns";
import Swal from "sweetalert2";
import { Button, Modal } from "flowbite-react";

const ListUser = () => {
  const [fetchStatus, setFetchStatus] = useState(true);
  const [jobs, setJobs] = useState([]);
  const [openModal, setOpenModal] = useState(false);
  const navigate = useNavigate();
  const [openModals, setOpenModals] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const pageSize = 4;
  const [activeModalUserId, setActiveModalUserId] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const token = Cookies.get("token");
        const response = await axios.get(
          `http://localhost:8080/api-blog-ngebruk/user?page=${currentPage}&pageSize=${pageSize}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        setJobs(response.data.users);
        setTotalPages(response.data.totalPages);
        window.dispatchEvent(new Event("userUpdated"));
      } catch (error) {
        console.error("Error fetching job list:", error);
      }
    };

    fetchData();
  }, [currentPage]); // Dependency on currentPage only

  const handleNext = () => {
    setCurrentPage((current) => current + 1);
  };

  const handlePrevious = () => {
    setCurrentPage((current) => Math.max(current - 1, 1));
  };

  const formatTanggal = (created_at) => {
    return format(new Date(created_at), " dd MMMM yyyy HH:mm:ss");
  };

  const handleEdit = (id_user) => {
    navigate(`/blogger/user/updateuser/${id_user}`);
  };

  const handleDelete = async (id_user) => {
    const token = Cookies.get("token");

    try {
      // Check if user is linked with a blogger
      const isUserLinked = await checkUserLinkedWithBlogger(id_user);

      if (isUserLinked) {
        Swal.fire(
          "Error",
          "User is linked with a blogger and cannot be deleted.",
          "error"
        );
        return;
      }

      // If not linked, proceed with deletion
      await axios.delete(
        `http://localhost:8080/api-blog-ngebruk/user/${id_user}`,
        {
          headers: {
            Authorization: `${token}`,
          },
        }
      );
      setActiveModalUserId(null);
      Swal.fire("Success", "User successfully deleted", "success");

      // Refresh the user list
      await refreshUserList();
    } catch (error) {
      console.error("Error deleting user:", error);
      Swal.fire("Error", "Failed to delete user", "error");
    }
  };

  const refreshUserList = async () => {
    try {
      const token = Cookies.get("token");
      const response = await axios.get(
        `http://localhost:8080/api-blog-ngebruk/user?page=${currentPage}&pageSize=${pageSize}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setJobs(response.data.users);
      setTotalPages(response.data.totalPages);
    } catch (error) {
      console.error("Error fetching job list:", error);
    }
  };

  const checkUserLinkedWithBlogger = async (id_user) => {
    try {
      const response = await axios.get(
        `http://localhost:8080/api-blog-ngebruk/user/${id_user}/blogger`
      );
      const bloggers = response.data;
      return bloggers.length > 0;
    } catch (error) {
      console.error("Error checking user-blogger link:", error);
      return false;
    }
  };

  const handleOpenModal = (id_user) => {
    setActiveModalUserId(id_user);
  };

  const handleCloseModal = () => {
    setActiveModalUserId(null);
  };

  return (
    <>
      <div className="">
        <div className="p-4 border-2 rounded-lg bg-gray-100 mt-20 font-serif">
          <h1 className="text-center font-serif text-4xl">Data Pengguna</h1>
          <div className="mt-2">
            <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
              <table className="table rounded-lg overflow-hidden w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400 ">
                <thead className="text-xs text-gray-700 uppercase bg-blue-400 dark:bg-gray-700 dark:text-gray-400">
                  <tr>
                    <th scope="col" className="px-6 py-3 text-white">
                      No
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Username
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Nama Pengguna
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Email
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Password
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Role
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Created At
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Updated At
                    </th>
                    <th scope="col" className="px-6 py-3 text-white">
                      Action
                    </th>
                  </tr>
                </thead>

                {jobs.map((job, numb) => (
                  <tbody key={job.id_user}>
                    <tr className="odd:bg-white odd:dark:bg-gray-900 even:bg-gray-50 even:dark:bg-gray-800 border-b dark:border-gray-700">
                      <th
                        scope="row"
                        className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                        key={numb}
                      >
                        {(currentPage - 1) * pageSize + numb + 1}
                      </th>
                      <td className="px-6 py-4 font-bold">
                        <h1>{job.username}</h1>
                      </td>
                      <td className="px-6 py-4 ">
                        <h1>{job.nama}</h1>
                      </td>
                      <td className="px-6 py-4 ">
                        <h1>{job.email}</h1>
                      </td>
                      <td className="px-6 py-4 ">
                        <h1>{job.password}</h1>
                      </td>
                      <td className="px-6 py-4 ">
                        <h1>{job.role}</h1>
                      </td>
                      <td className="px-6 py-4">
                        {formatTanggal(job.created_at)}
                      </td>
                      <td className="px-6 py-4">
                        {formatTanggal(job.UpdatedAt)}
                      </td>
                      <td className="px-6 py-4 space-x-4 space-x-reverse ">
                        <button
                          type="button"
                          value={job.id_user}
                          onClick={() => handleEdit(job.id_user)}
                          className=" text-white bg-yellow-300 border border-yellow-300 focus:outline-none hover:bg-yellow-100 focus:ring-4 focus:ring-yellow-200 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-yellow-800 dark:text-yellow dark:border-yellow-600 dark:hover:bg-yellow-700 dark:hover:border-yellow-600 dark:focus:ring-yellow-700"
                        >
                          Edit
                        </button>
                        <Button
                          color="failure"
                          onClick={() => handleOpenModal(job.id_user)}
                        >
                          Delete
                        </Button>
                        <Modal
                          show={activeModalUserId === job.id_user}
                          onDelete={() => handleDelete(job.id_user)}
                          onClose={handleCloseModal}
                        >
                          <Modal.Header>Warning</Modal.Header>
                          <Modal.Body>
                            <div className="space-y-6">
                              <p className="text-base leading-relaxed text-gray-500 dark:text-gray-400">
                                Apakah anda yakin mau dihapus ?
                              </p>
                            </div>
                          </Modal.Body>
                          <Modal.Footer>
                            <Button
                              color="failure"
                              onClick={() =>
                                handleDelete(job.id_user) && setOpenModal(false)
                              }
                            >
                              I accept
                            </Button>
                            <Button
                              color="gray"
                              onClick={() => handleCloseModal(numb)}
                            >
                              Decline
                            </Button>
                          </Modal.Footer>
                        </Modal>
                      </td>
                    </tr>
                  </tbody>
                ))}
              </table>
              <div className="p-4">
                <button
                  onClick={() => navigate("/blogger/user/createuser")}
                  type="button"
                  className=" focus:outline-none text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-red-700 dark:focus:ring-red-900"
                >
                  Tambah Pengguna
                </button>
              </div>
            </div>
          </div>
          <div className="flex justify-center mt-4">
            <button
              onClick={handlePrevious}
              disabled={currentPage === 1}
              className="flex items-center justify-center px-4 h-10 me-3 text-base font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
            >
              <svg
                className="w-3.5 h-3.5 me-2 rtl:rotate-180"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 14 10"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13 5H1m0 0 4 4M1 5l4-4"
                />
              </svg>
              Previous
            </button>
            <span className=" font-serif m-2 mr-5">Page {currentPage}</span>
            <button
              onClick={handleNext}
              disabled={currentPage >= totalPages}
              className="flex items-center justify-center px-4 h-10 text-base font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
            >
              Next
              <svg
                className="w-3.5 h-3.5 ms-2 rtl:rotate-180"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 14 10"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M1 5h12m0 0L9 1m4 4L9 9"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </>
  );
};
export default ListUser;
