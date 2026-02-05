import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Navbar from "./components/Navigation/Navbar";
import HomePage from "./components/HomePages";
import Blogger_Detail from "./components/Blogger_Detail";
import Login from "./components/Login";
import Sidebar from "./components/Navigation/Sidebar";
import AdminDashboard from "./components/AdminDashboard";
import AdminListBlogger from "./components/Admin/AdminListBlogger";
import CreateBloggerAdmin from "./components/Admin/FormBloggerAdmin";
import UpdateBloggerAdmin from "./components/Admin/FormBloggerAdmin";
import ListBloggerByUser from "./components/ListBloggerByUser";
import CreateBloggerUser from "./components/FormBloggerUser";
import UpdateBloggerUser from "./components/FormBloggerUser";
import ListKategori from "./components/Admin/ListKategori";
import CreateCategory from "./components/Admin/FormKategori";
import UpdateCategory from "./components/Admin/FormKategori";
import ListUser from "./components/Admin/ListUser";
import CreateUser from "./components/Admin/FormUserAdmin";
import UpdateUser from "./components/Admin/FormUserAdmin";
import UserById from "./components/UserById";
import FormUserWritter from "./components/FormUserWritter";
import ChangePassword from "./components/Auth/ChangePassword";
import PrivateRoute from "./components/Auth/PrivateRoute";
import ProfileDesa from "./components/ProfileDesa";
import Footer from "./components/Navigation/Footer";
import NotFound from "./components/NotFound";

function App() {
  return (
    <>
      <Router>
        <Routes>
          <Route
            path="/"
            element={
              <>
                <Navbar />
                <HomePage />
                <Footer />
              </>
            }
          />
          <Route
            path="/blogger_detail/:id"
            element={
              <>
                <Navbar />
                <Blogger_Detail />
                <Footer />
              </>
            }
          />
          <Route
            path="/profile-desa"
            element={
              <>
                <Navbar />
                <ProfileDesa />
                <Footer />
              </>
            }
          />
          <Route path="/blogger/loginakun" element={<Login />} />
          <Route
            path="/blogger/dashboard"
            element={
              <PrivateRoute allowedRoles={["admin", "writter"]}>
                <>
                  <Sidebar />
                  <AdminDashboard />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/listdashboardadmin"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <AdminListBlogger />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/create_blogger"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <CreateBloggerAdmin />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/update_blogger/:id_blogger"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <UpdateBloggerAdmin />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/createblogger"
            element={
              <PrivateRoute allowedRoles={["writter"]}>
                <>
                  <Sidebar />
                  <CreateBloggerUser />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/updateblogger/:id_blogger"
            element={
              <PrivateRoute allowedRoles={["writter"]}>
                <>
                  <Sidebar />
                  <UpdateBloggerUser />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/kategori"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <ListKategori />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/createkategori"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <CreateCategory />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/updatekategori/:id_kategori"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <UpdateCategory />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/all_user"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <ListUser />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/user/createuser"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <CreateUser />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/user/updateuser/:id_user"
            element={
              <PrivateRoute allowedRoles={["admin"]}>
                <>
                  <Sidebar />
                  <UpdateUser />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/listbloggeruser"
            element={
              <PrivateRoute allowedRoles={["writter"]}>
                <>
                  <Sidebar />
                  <ListBloggerByUser />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/userdetail"
            element={
              <PrivateRoute allowedRoles={["writter"]}>
                <>
                  <Sidebar /> <UserById />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/userdetail/updateuser/:id_user"
            element={
              <PrivateRoute allowedRoles={["writter"]}>
                <>
                  <Sidebar />
                  <FormUserWritter />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="/blogger/userdetail/change-password"
            element={
              <PrivateRoute allowedRoles={["writter"]}>
                <>
                  <Sidebar />
                  <ChangePassword />
                </>
              </PrivateRoute>
            }
          />
          <Route
            path="*"
            element={
              <>
                <Navbar />
                <NotFound />
              </>
            }
          />
        </Routes>
      </Router>
    </>
  );
}

export default App;
