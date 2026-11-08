/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface DsPatientsCategoriesFeedSerializer {
  age?: number;
  created_at?: string;
  creator_id?: number;
  description?: string;
  id?: number;
  image_url?: string;
  is_liked?: number;
  likes?: number;
  name?: string;
  updated_at?: string;
  video_url?: string;
  weight?: number;
}

export interface DsPatientsCategoriesSerializer {
  age?: number;
  created_at?: string;
  creator_id?: number;
  description?: string;
  id?: number;
  image_url?: string;
  likes?: number;
  name?: string;
  updated_at?: string;
  video_url?: string;
  weight?: number;
}

export interface DsPatientsCategoriesTileSerializer {
  age?: number;
  created_at?: string;
  creator_id?: number;
  description?: string;
  id?: number;
  image_url?: string;
  is_creator?: number;
  likes?: number;
  name?: string;
  updated_at?: string;
  video_url?: string;
  weight?: number;
}

export interface DsUserSerializer {
  id?: number;
  username?: string;
}

export interface HandlerLikeRequest {
  like?: number;
}

export interface HandlerLikeResponse {
  id?: number;
  likes?: number;
}

export interface HandlerLoginRequest {
  password?: string;
  username?: string;
}

export interface HandlerLoginResponse {
  access_token?: string;
  expires_in?: number;
  token_type?: string;
}

export interface HandlerPublishRequest {
  age?: number;
  description?: string;
  weight?: number;
}

export interface HandlerRegisterRequest {
  password?: string;
  username?: string;
}

import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  ResponseType,
} from "axios";
import axios from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams
  extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown>
  extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export type ContentType =
  | "application/json"
  | "application/vnd.api+json"
  | "multipart/form-data"
  | "application/x-www-form-urlencoded"
  | "text/plain";

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({
    securityWorker,
    secure,
    format,
    ...axiosConfig
  }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({
      ...axiosConfig,
      baseURL: axiosConfig.baseURL || "",
    });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected mergeRequestParams(
    params1: AxiosRequestConfig,
    params2?: AxiosRequestConfig,
  ): AxiosRequestConfig {
    const method = params1.method || (params2 && params2.method);

    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...((method &&
          this.instance.defaults.headers[
            method.toLowerCase() as keyof HeadersDefaults
          ]) ||
          {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected stringifyFormItem(formItem: unknown) {
    if (typeof formItem === "object" && formItem !== null) {
      return JSON.stringify(formItem);
    } else {
      return `${formItem}`;
    }
  }

  protected createFormData(input: Record<string, unknown>): FormData {
    if (input instanceof FormData) {
      return input;
    }
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      const propertyContent: any[] =
        property instanceof Array ? property : [property];

      for (const formItem of propertyContent) {
        const isFileType = formItem instanceof Blob || formItem instanceof File;
        formData.append(
          key,
          isFileType ? formItem : this.stringifyFormItem(formItem),
        );
      }

      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = format || this.format || undefined;

    if (
      type === "multipart/form-data" &&
      body &&
      body !== null &&
      typeof body === "object"
    ) {
      body = this.createFormData(body as Record<string, unknown>);
    }

    if (
      type === "text/plain" &&
      body &&
      body !== null &&
      typeof body !== "string"
    ) {
      body = JSON.stringify(body);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(requestParams.headers || {}),
        ...(type ? { "Content-Type": type } : {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title Body Mass Index API
 * @version 1.0
 * @contact
 *
 * Веб-сервис услуг. Без авторизации доступно только чтение данных
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  patientsCategories = {
    /**
     * @description Список опубликованных услуг с признаком is_creator; фильтр по возрасту ?age=N
     *
     * @tags patients_categories
     * @name PatientsCategoriesList
     * @summary Плитка опубликованных услуг
     * @request GET:/patients_categories
     */
    patientsCategoriesList: (
      query?: {
        /** возраст для фильтрации */
        age?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<DsPatientsCategoriesTileSerializer[], string>({
        path: `/patients_categories`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * @description Создаёт черновик с файлами изображения и видео; создатель берётся из JWT
     *
     * @tags patients_categories
     * @name PatientsCategoriesCreate
     * @summary Создать черновик услуги
     * @request POST:/patients_categories
     * @secure
     */
    patientsCategoriesCreate: (
      data: {
        /** название услуги */
        name: string;
        /**
         * изображение
         * @format binary
         */
        image?: File;
        /**
         * короткое видео
         * @format binary
         */
        video?: File;
      },
      params: RequestParams = {},
    ) =>
      this.request<DsPatientsCategoriesSerializer, string>({
        path: `/patients_categories`,
        method: "POST",
        body: data,
        secure: true,
        type: "multipart/form-data",
        format: "json",
        ...params,
      }),

    /**
     * @description Черновик пользователя из JWT, id не указывается (не более 1 записи)
     *
     * @tags patients_categories
     * @name DraftList
     * @summary Получить черновик текущего пользователя
     * @request GET:/patients_categories/draft
     * @secure
     */
    draftList: (params: RequestParams = {}) =>
      this.request<DsPatientsCategoriesSerializer, string>({
        path: `/patients_categories/draft`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Первая услуга или следующая после id при ?id=X&next=true; признак is_liked
     *
     * @tags patients_categories
     * @name FeedList
     * @summary Лента опубликованных услуг
     * @request GET:/patients_categories/feed
     */
    feedList: (
      query?: {
        /** id услуги */
        id?: number;
        /** получить следующую услугу */
        next?: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<DsPatientsCategoriesFeedSerializer, string>({
        path: `/patients_categories/feed`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * @description Soft delete: статус меняется на deleted; только услуги этого пользователя
     *
     * @tags patients_categories
     * @name PatientsCategoriesDelete
     * @summary Удалить услугу
     * @request DELETE:/patients_categories/{id}
     * @secure
     */
    patientsCategoriesDelete: (id: number, params: RequestParams = {}) =>
      this.request<string, string>({
        path: `/patients_categories/${id}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Поле like: 1 ставит лайк, 0 снимает; лайк от текущего пользователя
     *
     * @tags patients_categories
     * @name LikeCreate
     * @summary Поставить/снять лайк
     * @request POST:/patients_categories/{id}/like
     * @secure
     */
    likeCreate: (
      id: number,
      request: HandlerLikeRequest,
      params: RequestParams = {},
    ) =>
      this.request<HandlerLikeResponse, string>({
        path: `/patients_categories/${id}/like`,
        method: "POST",
        body: request,
        secure: true,
        type: "application/json",
        format: "json",
        ...params,
      }),

    /**
     * @description Заполняет description, age, weight и переводит статус draft -> published; только создатель услуги
     *
     * @tags patients_categories
     * @name PublishUpdate
     * @summary Опубликовать черновик услуги
     * @request PUT:/patients_categories/{id}/publish
     * @secure
     */
    publishUpdate: (
      id: number,
      request: HandlerPublishRequest,
      params: RequestParams = {},
    ) =>
      this.request<DsPatientsCategoriesSerializer, string>({
        path: `/patients_categories/${id}/publish`,
        method: "PUT",
        body: request,
        secure: true,
        type: "application/json",
        format: "json",
        ...params,
      }),
  };
  users = {
    /**
     * @description Проверяет логин и пароль, возвращает JWT
     *
     * @tags users
     * @name AuthCreate
     * @summary Аутентификация
     * @request POST:/users/auth
     */
    authCreate: (request: HandlerLoginRequest, params: RequestParams = {}) =>
      this.request<HandlerLoginResponse, string>({
        path: `/users/auth`,
        method: "POST",
        body: request,
        type: "application/json",
        format: "json",
        ...params,
      }),

    /**
     * @description Кладёт текущий токен в чёрный список Redis до истечения его срока
     *
     * @tags users
     * @name LogoutCreate
     * @summary Деавторизация
     * @request POST:/users/logout
     * @secure
     */
    logoutCreate: (params: RequestParams = {}) =>
      this.request<string, string>({
        path: `/users/logout`,
        method: "POST",
        secure: true,
        ...params,
      }),

    /**
     * @description Создаёт нового пользователя, пароль сохраняется в виде sha1-хэша
     *
     * @tags users
     * @name RegisterCreate
     * @summary Регистрация пользователя
     * @request POST:/users/register
     */
    registerCreate: (
      request: HandlerRegisterRequest,
      params: RequestParams = {},
    ) =>
      this.request<DsUserSerializer, string>({
        path: `/users/register`,
        method: "POST",
        body: request,
        type: "application/json",
        format: "json",
        ...params,
      }),
  };
}
