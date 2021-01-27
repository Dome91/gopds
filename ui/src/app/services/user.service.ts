import {Injectable} from '@angular/core';
import {User} from "../models/user";
import {Observable} from "rxjs";
import {HttpClient} from "@angular/common/http";
import {map} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor(private http: HttpClient) {
  }

  create(user: User): Observable<null> {
    return this.http.post<null>('/api/v1/users', {
      username: user.username,
      password: user.password,
      role: user.role,
    })
  }

  getAll(): Observable<User[]> {
    return this.http.get<GetUsersResponse>('/api/v1/users')
      .pipe(map(response => response.users));
  }

  delete(user: User): Observable<null> {
    return this.http.delete<null>(`/api/v1/users/${user.username}`)
  }
}

export class GetUsersResponse {
  constructor(public users: User[]) {
  }
}
