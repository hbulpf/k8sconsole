import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import {HttpClientModule} from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ServiceComponent } from './service/service.component';
import { NavbarComponent } from './navbar/navbar.component';
import { SearchComponent } from './search/search.component';
import { HomeComponent } from './home/home.component';
import { FooterComponent } from './footer/footer.component';
import { NamespaceComponent } from './namespace/namespace.component';
import { PodComponent } from './pod/pod.component';

@NgModule({
  declarations: [
    AppComponent,
    ServiceComponent,
    NavbarComponent,
    SearchComponent,
    HomeComponent,
    FooterComponent,
    NamespaceComponent,
    PodComponent
  ],
  imports: [
    HttpClientModule,
    BrowserModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
