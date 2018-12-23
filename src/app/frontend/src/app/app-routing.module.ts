import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {HomeComponent} from './home/home.component';
import {ServiceComponent} from './service/service.component';
import {NamespaceComponent} from './namespace/namespace.component';
import {PodComponent} from './pod/pod.component';

const routes: Routes = [
  {path: '', component: HomeComponent},
  {path: 'service', component: ServiceComponent},
  {path: 'namespace', component: NamespaceComponent},
  {path: 'pod', component: PodComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
