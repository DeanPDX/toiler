/**
 * A simple service to store session information.
 */
export class SessionService {
    static readonly tokenKey = 'auth-token';
    static readonly baseURL = '/api';

    /**
     * Get auth token. Will return empty string if not set.
     */
    public static getAuthToken(): string {
        let token = sessionStorage.getItem(this.tokenKey);
        if (token == null) {
            token = '';
        }
        return token;
    }

    /**
     * Set the auth token in session storage.
     * @param token The auth token you wish to set
     */
    public static setAuthToken(token: string) {
        sessionStorage.setItem(this.tokenKey, token);
    }

    /**
     * Returns true if user is authenticated.
     */
    public static isAuthenticated(): boolean {
        let token = sessionStorage.getItem(this.tokenKey);
        return token != null;
    }
}