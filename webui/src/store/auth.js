import { reactive } from 'vue';

export const authState = reactive({
    isLoggedIn: false, // It Tracks if the user is logged in
});

export function checkLoginStatus() {
    return !!localStorage.getItem('authToken');
}

export function logIn(token, id, username) {
    localStorage.setItem('authToken', token);
    localStorage.setItem('id', id);
    localStorage.setItem('username', username);
    authState.isLoggedIn = true;
}

export function logOut() {
    localStorage.removeItem('authToken');
    localStorage.removeItem('id');
    localStorage.removeItem('username');
    authState.isLoggedIn = false;
}

export function getToken() {
    return localStorage.getItem('authToken');

}

export function getId() {
    return localStorage.getItem('id');
}

export function getUsername() {
    return localStorage.getItem('username');
}
