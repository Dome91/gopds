import {TestBed} from '@angular/core/testing';

import {UserService} from './user.service';
import {Role, User} from "../models/user";
import {HttpClientTestingModule, HttpTestingController} from "@angular/common/http/testing";

describe('UserService', () => {
  let service: UserService;
  let client: HttpTestingController

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule]
    });

    service = TestBed.inject(UserService);
    client = TestBed.inject(HttpTestingController);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should create user', () => {
    const user = new User('user1', 'password1', Role.ADMIN);
    service.create(user).subscribe();

    const request = client.expectOne('/api/v1/users');
    expect(request.request.method).toBe('POST');
    expect(request.request.body).toEqual({username: 'user1', password: 'password1', role: Role.ADMIN});
    request.flush(null)
  });

  it('should fetch users', () => {
    const user1 = new User('user1', 'password1', Role.ADMIN);
    const user2 = new User('user2', 'password2', Role.USER);

    service.getAll().subscribe(
      response => {
        expect(response).toContain(user1);
        expect(response).toContain(user2);
      });

    const request = client.expectOne('/api/v1/users');
    expect(request.request.method).toBe('GET');
    request.flush({users: [user1, user2]});
  });

  it('should delete user', () => {
    const user = new User('user1', 'password1', Role.ADMIN);
    service.delete(user).subscribe();
    const request = client.expectOne('/api/v1/users/user1');
    expect(request.request.method).toBe('DELETE');
    request.flush(null);
  });
});
