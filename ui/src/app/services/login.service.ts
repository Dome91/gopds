import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class LoginService {

  constructor(private client: HttpClient) {
  }

  login(username: string, password: string): Observable<null> {
    return this.client.post<null>('/api/v1/login', {
      username,
      password
    });
  }

  isLoggedIn(): Observable<null> {
    return this.client.get<null>('/api/v1/login');
  }

  logout(): Observable<null> {
    return this.client.put<null>('/api/v1/logout', null);
  }
}
