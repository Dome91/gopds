import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {LoginComponent} from "./components/login/login.component";
import {CatalogComponent} from "./components/catalog/catalog.component";
import {SettingsComponent} from "./components/settings/settings.component";
import {CreateSourceComponent} from "./components/settings/settings-sources/create-source/create-source.component";
import {CreateUserComponent} from "./components/settings/settings-users/create-user/create-user.component";

const routes: Routes = [
  {path: 'catalog', component: CatalogComponent},
  {path: 'settings', component: SettingsComponent},
  {path: 'create-source', component: CreateSourceComponent},
  {path: 'create-user', component: CreateUserComponent},
  {path: '', component: LoginComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
