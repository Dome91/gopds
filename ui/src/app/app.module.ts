import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {LoginComponent} from './components/login/login.component';
import {NgbModule} from '@ng-bootstrap/ng-bootstrap';
import {ToastrModule} from "ngx-toastr";
import {HTTP_INTERCEPTORS, HttpClientModule} from "@angular/common/http";
import {FormsModule} from "@angular/forms";
import {CatalogComponent} from './components/catalog/catalog.component';
import {HttpErrorInterceptor} from "./interceptors/http-error.interceptor";
import { SettingsComponent } from './components/settings/settings.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { SettingsSourcesComponent } from './components/settings/settings-sources/settings-sources.component';
import { SettingsUsersComponent } from './components/settings/settings-users/settings-users.component';
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import { CreateSourceComponent } from './components/settings/settings-sources/create-source/create-source.component';
import { CreateUserComponent } from './components/settings/settings-users/create-user/create-user.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    CatalogComponent,
    SettingsComponent,
    NavbarComponent,
    SettingsSourcesComponent,
    SettingsUsersComponent,
    CreateSourceComponent,
    CreateUserComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule,
    NgbModule,
    ToastrModule.forRoot({
      positionClass: 'toast-bottom-right',
    }),
    FontAwesomeModule,
  ],
  providers: [{provide: HTTP_INTERCEPTORS, useClass: HttpErrorInterceptor, multi: true}],
  bootstrap: [AppComponent]
})
export class AppModule { }
