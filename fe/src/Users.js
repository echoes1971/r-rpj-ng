import React, { useContext, useEffect, useState } from "react";
import api from "./axios";
import { app_cfg } from "./app.cfg";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";

function Users() {
  const { t } = useTranslation();
  const [users, setUsers] = useState([]);
  const [query, setQuery] = useState("");
  const [editingUser, setEditingUser] = useState(null); // utente in modifica
  const { dark, themeClass } = useContext(ThemeContext);
  const endpoint = app_cfg.endpoint;
  console.log("Using endpoint:", endpoint);

    // Carica tutti gli utenti all'inizio
  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async (search = "") => {
    const token = localStorage.getItem("token");
    try {
      const res = await api.get(
        search
          ? endpoint + `/users?search=${encodeURIComponent(search)}`
          : endpoint + "/users",
        { headers: { Authorization: `Bearer ${token}` }, }
      );
      setUsers(res.data || []); // supponendo che l'API restituisca un array
    } catch (err) {
      console.log("Error loading users.");
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchUsers(query);
  };

  const handleEditClick = (user) => {
    // alert("Editing user="+user.Login);
    setEditingUser({ ...user }); // copia per non modificare direttamente
  };

  const handleEditChange = (e) => {
    const { name, value } = e.target;
    // alert(name+"="+value);
    setEditingUser((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    const token = localStorage.getItem("token");
    try {
      if (!editingUser.ID) {
        // Nuovo utente
        await api.post(
          endpoint + `/users`,
          editingUser,
          { headers: { Authorization: `Bearer ${token}` } }
        );
        setEditingUser(null);
        fetchUsers();
        return;
      }
      // Utente esistente
      await api.put(
        endpoint + `/users/${editingUser.ID}`,
        editingUser,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setEditingUser(null);
      fetchUsers();
    } catch (err) {
      alert("Errore salvataggio utente");
    }
  };

  const handleDelete = async () => {
    if (window.confirm("Are you sure to delete this user?")) {
      const token = localStorage.getItem("token");
      await api.delete(endpoint + `/users/${editingUser.ID}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setEditingUser(null);
      fetchUsers();
    }
  };

  return (
    <div className={`container mt-3 ${themeClass}`}>
      <h2 className={dark ? "text-light" : "text-dark"}>{t("users.users")}</h2>

      {/* Form di ricerca */}
      {!editingUser && (
        <form className="d-flex mb-3" onSubmit={handleSearch}>
          <input
            type="text"
            className={`form-control me-2 ${dark ? "bg-secondary text-light" : ""}`}
            placeholder={t("common.search")}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <button className="btn btn-primary">{t("common.search")}</button>
        </form>
      )}

      {!editingUser && (
        <button
          className="btn btn-success mb-3"
          onClick={() => setEditingUser({ ID: "", Login: "", Fullname: "", GroupID: "" })}
        >
          {t("users.newUser")}
        </button>
      )}

      {!editingUser && users.length > 0 && (
        <table 
        className={`table ${dark ? "table-dark" : "table-striped"} table-hover`}
        >
          <thead>
            <tr>
              <th className="d-none d-md-table-cell">ID</th>
              <th>Login</th>
              <th>Fullname</th>
              <th className="d-none d-md-table-cell">Group</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map(u => (
              <tr key={u.ID}>
                <td className="d-none d-md-table-cell">{u.ID}</td>
                <td>{u.Login}</td>
                <td>{u.Fullname}</td>
                <td className="d-none d-md-table-cell">{u.GroupID}</td>
                <td>
                  <button
                    className="btn btn-sm btn-warning"
                    onClick={() => handleEditClick(u)}
                  >
                    {t("common.edit")}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* Form di editing */}
      {editingUser && (
        <div className={`card p-3 mt-3 ${dark ? "bg-dark text-light" : "bg-light text-dark" }`}>
          <h4>{editingUser.ID>'' ? t("common.edit") : t("common.create")} {t("users.user")}</h4>
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="ID"
            title="ID"
            value={editingUser.ID}
            readOnly
            onChange={handleEditChange}
          />
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="Login"
            title="Login"
            value={editingUser.Login}
            {...editingUser.ID ? "disabled" : null}
            onChange={handleEditChange}
          />
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="Fullname"
            title="Fullname"
            value={editingUser.Fullname}
            onChange={handleEditChange}
          />
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="GroupID"
            title="GroupID"
            value={editingUser.GroupID}
            readOnly
            onChange={handleEditChange}
          />
          <div>
            <button className="btn btn-success me-2" onClick={handleSave}>
              { t("common.save") }
            </button>
            <button
              className="btn btn-secondary me-4"
              onClick={() => setEditingUser(null)}
            >
              { t("common.cancel") }
            </button>
            {editingUser.ID>"" ?
                      <button
                        className="btn btn-danger"
                        onClick={handleDelete}
                      >
                        { t("common.delete") }
                      </button>
          : null}
          </div>
        </div>
      )}
    </div>
  );
}

export default Users;
